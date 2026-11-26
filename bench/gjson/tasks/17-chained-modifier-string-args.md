kind: neutral
Issue #143: chaining modifiers that take a string argument stops after the first one. With a custom modifier `push` that appends its argument to an array,

    @push:"2"|@push:"3"|@push:{"a":"b","c":["e","f"]}|@push:true|@push:10.23

on `[]` should return `["2","3",{"a":"b","c":["e","f"]},true,10.23]`. Today only the first modifier runs, because its quoted argument swallows the rest of the path. An argument is a JSON value (string, object, array, number, literal) and ends where that value ends. Please fix and add a test.
