package wolframfns

// Function is a single mathematical function from functions.wolfram.com.
type Function struct {
	Rank     int    `json:"rank"     csv:"rank"     tsv:"rank"`
	Category string `json:"category" csv:"category" tsv:"category"`
	Name     string `json:"name"     csv:"name"     tsv:"name"`
	URL      string `json:"url"      csv:"url"      tsv:"url"`
}

// Category is a grouping of functions with a count.
type Category struct {
	Rank  int    `json:"rank"  csv:"rank"  tsv:"rank"`
	Name  string `json:"name"  csv:"name"  tsv:"name"`
	Count int    `json:"count" csv:"count" tsv:"count"`
}
