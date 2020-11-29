# -add-options json=omitempty

gen:
	ls *.go | xargs -I {} -n 1 gomodifytags -file={} -w -all -add-tags json -transform camelcase --skip-unexported
