package service

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// getAllCatalogItemsHandler returns a fake list of catalog items
func getAllCatalogItemsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		catalog := make([]catalogItem, 2)
		formatter.JSON(w, http.StatusOK, catalog)
	}
}

// getCatalogItemDetailsHandler returns a fake catalog item. The key takeaway here
// is that we're using a backing service to get fulfillment status for the individual
// item.
func getCatalogItemDetailsHandler(formatter *render.Render, serviceClient fulfillmentClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		sku := vars["sku"]
		status, err := serviceClient.getFulfillmentStatus(sku)
		if err == nil {
			formatter.JSON(w, http.StatusOK, catalogItem{
				ProductID:       1,
				SKU:             sku,
				Description:     "This is a fake product",
				Price:           1599, // $15.99
				ShipsWithin:     status.ShipsWithin,
				QuantityInStock: status.QuantityInStock,
			})
		} else {
			formatter.JSON(w, http.StatusInternalServerError, fmt.Sprintf("Fulfillment Client error: %s", err.Error()))
		}
	}
}