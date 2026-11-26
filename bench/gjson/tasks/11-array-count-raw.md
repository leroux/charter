kind: neutral
The count from `#` has an empty `Raw`: `gjson.Get("[1,2,3,4]", "#").Raw` is `""`, while `.String()` and `.Int()` are fine. Code that forwards `Raw` sends nothing. `Raw` should be `4`. Please fix it and add a test.
