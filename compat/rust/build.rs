extern crate cbindgen;

use std::env;

fn main() {
    cbindgen::Builder::new()
        .with_language(cbindgen::Language::C)
        .with_crate(env::var("CARGO_MANIFEST_DIR").expect("CARGO_MANIFEST_DIR is not set/valid"))
        .generate()
        .expect("unable to generate bindings")
        .write_to_file("compat.h");
}
