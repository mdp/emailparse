# emailparse

_Be aware: very rough first pass for a small project. Not really well packaged at this point_

CLI email parser to template.

## Example

```
cat mime_email.txt | emailparse "{{ datef \"2006.01.02\"}}: {{.Subject}} - {{ .Text }}" > email.txt
```

## License

MIT
