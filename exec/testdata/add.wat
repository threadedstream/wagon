(module
    (import "logger" "log" (func $log (param $i i32)))
    (func $add (param $x i32) (param $y i32) (result i32)
        (i32.add 
            (local.get $x)
            (local.get $y)
        )
    )
    (func $add_printed (param $x i32) (param $y i32) 
        (local $temp i32)
        (i32.add 
            (local.get $x)
            (local.get $y)
        )
        (local.set $temp)
        (call $log (local.get $temp))
    )
    (export "add" (func $add))
)