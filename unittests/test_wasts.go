package unittests

var i32_overflow_wast string = `(module
 (import "env" "require_auth" (func $require_auth (param i64)))
 (import "env" "eosio_assert" (func $eosio_assert (param i32 i32)))
  (table 0 anyfunc)
  (memory $0 1)
  (export "apply" (func $apply))
  (func $i32_trunc_s_f32 (param $0 f32) (result i32) (i32.trunc_s/f32 (get_local $0)))
  (func $i32_trunc_u_f32 (param $0 f32) (result i32) (i32.trunc_u/f32 (get_local $0)))
  (func $i32_trunc_s_f64 (param $0 f64) (result i32) (i32.trunc_s/f64 (get_local $0)))
  (func $i32_trunc_u_f64 (param $0 f64) (result i32) (i32.trunc_u/f64 (get_local $0)))
  (func $test (param $0 i32))
  (func $apply (param $0 i64)(param $1 i64)(param $2 i64)
   (call $test (call $%s (%s)))
))`

var i64_overflow_wast string = `(module
  (import "env" "require_auth" (func $require_auth (param i64)))
  (import "env" "eosio_assert" (func $eosio_assert (param i32 i32)))
   (table 0 anyfunc)
   (memory $0 1)
   (export "apply" (func $apply))
   (func $i64_trunc_s_f32 (param $0 f32) (result i64) (i64.trunc_s/f32 (get_local $0)))
   (func $i64_trunc_u_f32 (param $0 f32) (result i64) (i64.trunc_u/f32 (get_local $0)))
   (func $i64_trunc_s_f64 (param $0 f64) (result i64) (i64.trunc_s/f64 (get_local $0)))
   (func $i64_trunc_u_f64 (param $0 f64) (result i64) (i64.trunc_u/f64 (get_local $0)))
   (func $test (param $0 i64))
   (func $apply (param $0 i64)(param $1 i64)(param $2 i64)
    (call $test (call $%s (%s)))
))`

var aligned_ref_wast string = `(module
 (import "env" "sha256" (func $sha256 (param i32 i32 i32)))
 (table 0 anyfunc)
 (memory $0 32)
 (data (i32.const 4) "hello")
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (call $sha256
   (i32.const 4)
   (i32.const 5)
   (i32.const 16)
  )
 )
)`

var aligned_const_ref_wast string = `(module
 (import "env" "sha256" (func $sha256 (param i32 i32 i32)))
 (import "env" "assert_sha256" (func $assert_sha256 (param i32 i32 i32)))
 (table 0 anyfunc)
 (memory $0 32)
 (data (i32.const 4) "hello")
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (local $3 i32)
  (call $sha256
   (i32.const 4)
   (i32.const 5)
   (i32.const 16)
  )
  (call $assert_sha256
   (i32.const 4)
   (i32.const 5)
   (i32.const 16)
  )
 )
)`

var misaligned_ref_wast string = `(module
 (import "env" "sha256" (func $sha256 (param i32 i32 i32)))
 (table 0 anyfunc)
 (memory $0 32)
 (data (i32.const 4) "hello")
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (call $sha256
   (i32.const 4)
   (i32.const 5)
   (i32.const 5)
  )
 )
)`

var misaligned_const_ref_wast string = `(module
 (import "env" "sha256" (func $sha256 (param i32 i32 i32)))
 (import "env" "assert_sha256" (func $assert_sha256 (param i32 i32 i32)))
 (import "env" "memmove" (func $memmove (param i32 i32 i32) (result i32)))
 (table 0 anyfunc)
 (memory $0 32)
 (data (i32.const 4) "hello")
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (local $3 i32)
  (call $sha256
   (i32.const 4)
   (i32.const 5)
   (i32.const 16)
  )
  (set_local $3
   (call $memmove
    (i32.const 17)
    (i32.const 16)
    (i32.const 64) 
   )
  )
  (call $assert_sha256
   (i32.const 4)
   (i32.const 5)
   (i32.const 17)
  )
 )
)`

var entry_wast string = `(module
 (import "env" "require_auth" (func $require_auth (param i64)))
 (import "env" "eosio_assert" (func $eosio_assert (param i32 i32)))
 (import "env" "current_time" (func $current_time (result i64)))
 (table 0 anyfunc)
 (memory $0 1)
 (export "memory" (memory $0))
 (export "entry" (func $entry))
 (export "apply" (func $apply))
 (func $entry
  (block
   (i64.store offset=4
    (i32.const 0)
    (call $current_time)
   )
  )
 )
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (block
   (call $require_auth (i64.const 6121376101093867520))
   (call $eosio_assert
    (i64.eq
     (i64.load offset=4
      (i32.const 0)
     )
     (call $current_time)
    )
    (i32.const 0)
   )
  )
 )
 (start $entry)
)`

var entry_wast_2 string = `(module
 (import "env" "require_auth" (func $require_auth (param i64)))
 (import "env" "eosio_assert" (func $eosio_assert (param i32 i32)))
 (import "env" "current_time" (func $current_time (result i64)))
 (table 0 anyfunc)
 (memory $0 1)
 (export "memory" (memory $0))
 (export "apply" (func $apply))
 (start $entry)
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (block
   (call $require_auth (i64.const 6121376101093867520))
   (call $eosio_assert
    (i64.eq
     (i64.load offset=4
      (i32.const 0)
     )
     (call $current_time)
    )
    (i32.const 0)
   )
  )
 )
 (func $entry
  (block
   (i64.store offset=4
    (i32.const 0)
    (call $current_time)
   )
  )
 )
)`

var biggest_memory_wast string = `(module
 (import "env" "eosio_assert" (func $$eosio_assert (param i32 i32)))
 (import "env" "require_auth" (func $$require_auth (param i64)))
 (table 0 anyfunc)
 (memory $$0 %d)
 (export "memory" (memory $$0))
 (export "apply" (func $$apply))
 (func $$apply (param $$0 i64) (param $$1 i64) (param $$2 i64)
  (call $$require_auth (i64.const 4294504710842351616))
  (call $$eosio_assert
   (i32.eq
     (grow_memory (i32.const 1))
     (i32.const -1)
   )
   (i32.const 0)
  )
 )
)`

var too_big_memory_wast string = `(module
 (table 0 anyfunc)
 (memory $$0 %d)
 (export "memory" (memory $$0))
 (export "apply" (func $$apply))
 (func $$apply (param $$0 i64) (param $$1 i64) (param $$2 i64))
)`

var simple_no_memory_wast string = `(module
 (import "env" "require_auth" (func $require_auth (param i64)))
 (import "env" "memcpy" (func $memcpy (param i32 i32 i32) (result i32)))
 (table 0 anyfunc)
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
    (call $require_auth (i64.const 11323361180581363712))
    (drop
       (call $memcpy
          (i32.const 0)
          (i32.const 1024)
          (i32.const 1024)
       )
    )
 )
)`

var mutable_global_wast string = `(module
 (import "env" "require_auth" (func $require_auth (param i64)))
 (import "env" "eosio_assert" (func $eosio_assert (param i32 i32)))
 (table 0 anyfunc)
 (memory $0 1)
 (export "memory" (memory $0))
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
  (call $require_auth (i64.const 7235159549794234880))
  (if (i64.eq (get_local $2) (i64.const 0)) (then
    (set_global $g0 (i64.const 444))
    (return)
  ))
  (if (i64.eq (get_local $2) (i64.const 1)) (then
    (call $eosio_assert (i64.eq (get_global $g0) (i64.const 2)) (i32.const 0))
    (return)
  ))
  (call $eosio_assert (i32.const 0) (i32.const 0))
 )
 (global $g0 (mut i64) (i64.const 2))
)`

var valid_sparse_table string = `(module
 (table 1024 anyfunc)
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64))
 (elem (i32.const 0) $apply)
 (elem (i32.const 1022) $apply $apply)
)`

var too_big_table string = `(module
 (table 1025 anyfunc)
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64))
 (elem (i32.const 0) $apply)
 (elem (i32.const 1022) $apply $apply)
)`

var memory_init_borderline string = `(module
 (memory $0 16)
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64))
 (data (i32.const 65532) "sup!")
)`

var memory_init_toolong string = `(module
 (memory $0 16)
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64))
 (data (i32.const 65533) "sup!")
)`

var memory_init_negative string = `(module
 (memory $0 16)
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64))
 (data (i32.const -1) "sup!")
)`

var memory_table_import string = `(module
 (table  (import "foo" "table") 10 anyfunc)
 (memory (import "nom" "memory") 0)
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64))
)`

var table_checker_wast string = `(module
 (import "env" "require_auth" (func $require_auth (param i64)))
 (import "env" "eosio_assert" (func $assert (param i32 i32)))
 (type $SIG$vj (func (param i64)))
 (table 1024 anyfunc)
 (memory $0 1)
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
   (call $require_auth (i64.const 14547189746360123392))
   (call_indirect $SIG$vj
     (i64.shr_u
       (get_local $2)
       (i64.const 32)
     )
     (i32.wrap/i64
       (get_local $2)
     )
   )
 )
 (func $apple (type $SIG$vj) (param $0 i64)
   (call $assert
     (i64.eq
       (get_local $0)
       (i64.const 555)
     )
     (i32.const 0)
   )
 )
 (func $bannna (type $SIG$vj) (param $0 i64)
   (call $assert
     (i64.eq
       (get_local $0)
       (i64.const 7777)
     )
     (i32.const 0)
   )
 )
 (elem (i32.const 0) $apple)
 (elem (i32.const 1022) $apple $bannna)
)`

var table_checker_proper_syntax_wast string = `(module
 (import "env" "require_auth" (func $require_auth (param i64)))
 (import "env" "eosio_assert" (func $assert (param i32 i32)))
 (import "env" "printi" (func $printi (param i64)))
 (type $SIG$vj (func (param i64)))
 (table 1024 anyfunc)
 (memory $0 1)
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
   (call $require_auth (i64.const 14547189746360123392))
   (call_indirect (type $SIG$vj)
     (i64.shr_u
       (get_local $2)
       (i64.const 32)
     )
     (i32.wrap/i64
       (get_local $2)
     )
   )
 )
 (func $apple (type $SIG$vj) (param $0 i64)
   (call $assert
     (i64.eq
       (get_local $0)
       (i64.const 555)
     )
     (i32.const 0)
   )
 )
 (func $bannna (type $SIG$vj) (param $0 i64)
   (call $assert
     (i64.eq
       (get_local $0)
       (i64.const 7777)
     )
     (i32.const 0)
   )
 )
 (elem (i32.const 0) $apple)
 (elem (i32.const 1022) $apple $bannna)
)`

var table_checker_small_wast string = `(module
 (import "env" "require_auth" (func $require_auth (param i64)))
 (import "env" "eosio_assert" (func $assert (param i32 i32)))
 (import "env" "printi" (func $printi (param i64)))
 (type $SIG$vj (func (param i64)))
 (table 128 anyfunc)
 (memory $0 1)
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
   (call $require_auth (i64.const 14547189746360123392))
   (call_indirect (type $SIG$vj)
     (i64.shr_u
       (get_local $2)
       (i64.const 32)
     )
     (i32.wrap/i64
       (get_local $2)
     )
   )
 )
 (func $apple (type $SIG$vj) (param $0 i64)
   (call $assert
     (i64.eq
       (get_local $0)
       (i64.const 555)
     )
     (i32.const 0)
   )
 )
 (elem (i32.const 0) $apple)
)`

var global_protection_none_get_wast string = `(module
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
   (drop (get_global 0))
 )
)`

var global_protection_some_get_wast string = `(module
 (global i32 (i32.const -11))
 (global i32 (i32.const 56))
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
   (drop (get_global 1))
   (drop (get_global 2))
 )
)`

var global_protection_none_set_wast string = `(module
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
   (set_global 0 (get_local 1))
 )
)`

var global_protection_some_set_wast string = `(module
 (global i64 (i64.const -11))
 (global i64 (i64.const 56))
 (export "apply" (func $apply))
 (func $apply (param $0 i64) (param $1 i64) (param $2 i64)
   (set_global 2 (get_local 1))
 )
)`

var global_protection_okay_get_wasm = []byte{
	0x00, 'a', 's', 'm', 0x01, 0x00, 0x00, 0x00,
	0x01, 0x07, 0x01, 0x60, 0x03, 0x7e, 0x7e, 0x7e, 0x00, //type section containing a function as void(i64,i64,i64)
	0x03, 0x02, 0x01, 0x00, //a function

	0x06, 0x06, 0x01, 0x7f, 0x00, 0x41, 0x75, 0x0b, //global

	0x07, 0x09, 0x01, 0x05, 'a', 'p', 'p', 'l', 'y', 0x00, 0x00, //export function 0 as "apply"
	0x0a, 0x07, 0x01, //code section
	0x05, 0x00, //function body start with length 5; no locals
	0x23, 0x00, //get global 0
	0x1a, //drop
	0x0b, //end
}

var global_protection_none_get_wasm = []byte{
	0x00, 'a', 's', 'm', 0x01, 0x00, 0x00, 0x00,
	0x01, 0x07, 0x01, 0x60, 0x03, 0x7e, 0x7e, 0x7e, 0x00, //type section containing a function as void(i64,i64,i64)
	0x03, 0x02, 0x01, 0x00, //a function

	0x07, 0x09, 0x01, 0x05, 'a', 'p', 'p', 'l', 'y', 0x00, 0x00, //export function 0 as "apply"
	0x0a, 0x07, 0x01, //code section
	0x05, 0x00, //function body start with length 5; no locals
	0x23, 0x00, //get global 0
	0x1a, //drop
	0x0b, //end
}

var global_protection_some_get_wasm = []byte{
	0x00, 'a', 's', 'm', 0x01, 0x00, 0x00, 0x00,
	0x01, 0x07, 0x01, 0x60, 0x03, 0x7e, 0x7e, 0x7e, 0x00, //type section containing a function as void(i64,i64,i64)
	0x03, 0x02, 0x01, 0x00, //a function

	0x06, 0x06, 0x01, 0x7f, 0x00, 0x41, 0x75, 0x0b, //global

	0x07, 0x09, 0x01, 0x05, 'a', 'p', 'p', 'l', 'y', 0x00, 0x00, //export function 0 as "apply"
	0x0a, 0x07, 0x01, //code section
	0x05, 0x00, //function body start with length 5; no locals
	0x23, 0x01, //get global 1
	0x1a, //drop
	0x0b, //end
}

var global_protection_okay_set_wasm = []byte{
	0x00, 'a', 's', 'm', 0x01, 0x00, 0x00, 0x00,
	0x01, 0x07, 0x01, 0x60, 0x03, 0x7e, 0x7e, 0x7e, 0x00, //type section containing a function as void(i64,i64,i64)
	0x03, 0x02, 0x01, 0x00, //a function

	0x06, 0x06, 0x01, 0x7e, 0x01, 0x42, 0x75, 0x0b, //global.. this time with i64 & global mutablity

	0x07, 0x09, 0x01, 0x05, 'a', 'p', 'p', 'l', 'y', 0x00, 0x00, //export function 0 as "apply"
	0x0a, 0x08, 0x01, //code section
	0x06, 0x00, //function body start with length 6; no locals
	0x20, 0x00, //get local 0
	0x24, 0x00, //set global 0
	0x0b, //end
}

var global_protection_some_set_wasm = []byte{
	0x00, 'a', 's', 'm', 0x01, 0x00, 0x00, 0x00,
	0x01, 0x07, 0x01, 0x60, 0x03, 0x7e, 0x7e, 0x7e, 0x00, //type section containing a function as void(i64,i64,i64)
	0x03, 0x02, 0x01, 0x00, //a function

	0x06, 0x06, 0x01, 0x7e, 0x01, 0x42, 0x75, 0x0b, //global.. this time with i64 & global mutablity

	0x07, 0x09, 0x01, 0x05, 'a', 'p', 'p', 'l', 'y', 0x00, 0x00, //export function 0 as "apply"
	0x0a, 0x08, 0x01, //code section
	0x06, 0x00, //function body start with length 6; no locals
	0x20, 0x00, //get local 0
	0x24, 0x01, //set global 1
	0x0b, //end
}