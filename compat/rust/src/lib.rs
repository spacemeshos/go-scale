#![feature(vec_into_raw_parts)]

use parity_scale_codec::{Compact, Decode, Encode, Error};
use std::ffi::c_uchar;

#[derive(Encode, Decode)]
struct Struct {
    field1: Compact<u16>,
    field2: [u8; 3],
}

macro_rules! decode_encode {
	($input:ident, $output:ident, $($type:ty),*) => {
        $(<$type>::decode($input)?.encode_to(&mut $output);)*
    };
}

fn decode_encode(buf: &[u8]) -> Result<Vec<u8>, Error> {
    let input = &mut &buf[..];
    let mut output = vec![];
    decode_encode!(
        input,
        output,
        Compact<u8>,
        Compact<u16>,
        Compact<u32>,
        Compact<u64>,
        [u8; 8],
        bool,
        Option<Struct>,
        Struct,
        [Struct; 4],
        Vec<u8>,
        Vec<[u8; 32]>,
        Vec<Vec<u8>>,
        Vec<Struct>
    );
    Ok(output)
}

#[repr(C)]
pub struct Response {
    code: usize,
    ptr: *mut u8,
    len: usize,
}

/// # Safety
/// 
/// input must be a valid pointer.
#[no_mangle]
pub unsafe extern "C" fn round_trip(input: *const c_uchar, input_len: usize) -> *mut Response {
    let input = unsafe { std::slice::from_raw_parts(input, input_len) };
    let resp = Box::new(match decode_encode(input) {
        Ok(output) => {
            let (ptr, len, _) = output.into_raw_parts();
            Response {
                code: 0,
                ptr,
                len,
            }
        }
        Err(err) => {
            let s = err.to_string();
            let (ptr, len, _) = s.into_raw_parts();
            Response {
                code: 1,
                ptr,
                len,
            }
        }
    });
    Box::into_raw(resp)
}

#[no_mangle]
pub unsafe extern "C" fn free_response(response: *mut Response) {
    let response = Box::from_raw(response);
    match response.code {
        0 => {
            Vec::from_raw_parts(response.ptr, response.len, response.len);
        },
        1 => {
            String::from_raw_parts(response.ptr, response.len, response.len);
        },
        _ => {},
    } 
}
