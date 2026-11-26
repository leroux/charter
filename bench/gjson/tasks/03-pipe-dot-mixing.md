kind: neutral
I would like to mix pipes and dots freely in one path. On

    {"info":{"friends":[{"first":"Dale","last":"Murphy","age":44},{"first":"Roger","last":"Craig","age":68},{"first":"Jane","last":"Murphy","age":47}]}}

all of these should return `Murphy`:

    info.friends.#[first="Dale"].last
    info|friends.#[first="Dale"].last
    info|friends.#[first="Dale"]|last
    info|friends|#[first="Dale"]|last
    @ugly|info|friends|#[first="Dale"]|last
    @ugly|info.@ugly|friends|#[first="Dale"]|last
    @ugly.info|@ugly.friends|#[first="Dale"]|last

so a modifier may also follow a dot (`info.@ugly`), not only a pipe. Please add tests.
