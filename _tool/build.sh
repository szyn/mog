glide install -v
gox -osarch "linux/amd64 linux/arm darwin/amd64 windows/amd64" \
-parallel=4 -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}"
