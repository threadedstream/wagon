(module 
    (import "go_foo" "Print" 
    (func $print (param $x i32)))

    (import "go_foo" "DisplayHelloWorld"
    (func $display_hello_world ))

    (global $x1 i32 (i32.const 80))
    (global $x2 i32 (i32.const 20))

    (func $add
        (local $temp i32)
        (local.set $temp (
        i32.add 
            (global.get $x1)
            (global.get $x2)
        ))
        (call $print (local.get $temp))
        (call $display_hello_world)
    )

    (func $second_add 
        (call $add)
    )
    (export "add" (func $add))
)