# xorf

Simple tool for XOR encryption of files. Example of usage:

```bash
# Encrypt file text.txt using key files key1 and key2
$ xorf text.txt key1 key2 > encrypted1.txt

# Create new key, encrypt file encrypted1 by it, write key to file key3
$ xorf encrypted1.txt -k key3 > encrypted2.txt

# Decrypt file by all the keys
$ xorf encrypted2.txt key1 key2 key3 > decrypted.txt

# text.txt and decrypted.txt are the same
$ cmp text.txt decrypted.txt
```

`xorf` ends writing to stdout when one of file ends, so encrypted file will have the same size as the smallest file.

Please, don't pipe `xorf` output to the one of the source files. It can cause bugs.
