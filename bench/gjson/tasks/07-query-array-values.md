kind: neutral
Queries only work on arrays of objects. I would like to query an array of plain values too. On

    {"artists":[["Bob Dylan"],"John Lennon","Mick Jagger","Elton John","Michael Jackson","John Smith",true,123,456,false,null]}

- `a*.#[0="Bob Dylan"]#|#` should be `1` (the query path `0` reads into the nested array),
- `a*.#[0="Bob Dylan 2"]#|#` should be `0`,
- `a*.#[%"John*"]#|#` should be `2` (an empty query path means the element itself),
- `a*.#[_%"John*"]#|#` should be `0` (`_` is a key the strings do not have, not the element),
- `a*.#[="123"]#|#` should be `1`.

Please add it with tests.
