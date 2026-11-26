kind: neutral
when a file doesn't exist the tools print go's error ("open x: no such file or directory") and stop. make them behave like gnu: same error messages on stderr, keep going with the other files, right exit code
