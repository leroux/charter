kind: neutral
add an encoder too: Go values -> TOML text in the package, plus a cmd/toml-encoder that does toml-test's encoder side (tagged JSON on stdin, TOML on stdout, non-zero exit if it can't). output should decode back to the same values
