package horo

import "golang.org/x/net/context"

func ExampleBasic() {
	h := New()
	h.GET("/", func(c context.Context) error {
		return Text(c, 200, "ok")
	})
	h.ListenAndServe(":8080")
}
