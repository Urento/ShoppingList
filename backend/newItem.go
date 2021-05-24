package hello

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	Structs "github.com/urento/shoppinglist/structs"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	var shoppingList Structs.ShoppingList
	if err := json.NewDecoder(r.Body).Decode(&shoppingList); err != nil {
		fmt.Fprint(w, "Nichts")
	}
	fmt.Fprintf(w, "Hello, %s!", html.EscapeString(shoppingList.Title))
}
