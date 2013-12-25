var MPAjax = (function() {
    var MPAjax = {};
    
    function doAjax(data, success) {
        $.ajax({
            url: "/api",
            type: "POST",
            dataType: "json",
            data: data,
            
            error: function(jqXHR, textStatus, errorThrown) {
                console.log("MPAjax HTTP error:");
                console.log("  textStatus = " + textStatus);
                console.log("  errorThrown = " + errorThrown);
                alert("MPAjax error! Check console for more details.");
            },
            
            success: function(resp, textStatus, jqXHR) {
                if (resp.error) {
                    console.log("MPAjax server error:");
                    console.log("  Message: " + resp.error);
                }
                
                else {
                    success(resp.success);
                }
            },
        });
    }
    
    MPAjax.fetchMealList = function(destElement) {
        var params = {
            "command": "fetch-meal-list",
        };
        
        doAjax(params, function(results) {
            results = results || [];
            
            if (results.length == 0) {
                $(destElement).text("No meals to display.");
                return;
            }
            
            var table = $("<table>").appendTo($(destElement).empty());
            var headerRow = $("<tr>").appendTo(table);
            $("<th>").text("Name").appendTo(headerRow);
            $("<th>").text("Tags").appendTo(headerRow);
            $("<th>").text("Actions").appendTo(headerRow);
            
            var i, result, row, actions, favText;
            
            for (i = 0; i < results.length; i++) {
                result = results[i];
                row = $("<tr>").appendTo(table);
                
                $("<td>").appendTo(row).text(result.meal.name);
                $("<td>").appendTo(row).text((result.tags || []).join(", "));
                actions = $("<td>").appendTo(row);
                
                $("<button>").appendTo(actions).addClass("action-button").text("Open recipe").click(function(event) {
                    event.preventDefault();
                    location.href = result.meal.recipe;
                });
                
                favText = result.meal.favourite ? "Unfavourite" : "Favourite";
                $("<button>").appendTo(actions).addClass("action-button").text(favText).click(function(event) {
                    event.preventDefault();
                    MPAjax.toggleFavourite(result.meal.id, this);
                });
                
                $("<button>").appendTo(actions).addClass("action-button").text("Edit").click(function(event) {
                    event.preventDefault();
                    location.href = "/meals/" + result.meal.id + "/edit";
                });
                
                $("<button>").appendTo(actions).addClass("action-button").text("Delete").click(function(event) {
                    event.preventDefault();
                    MPAjax.deleteMeal(result.meal.id, row[0]);
                });
            }
        });
    };
    
    MPAjax.toggleFavourite = function(mealID, favButton) {
        var params = {
            "command": "toggle-favourite",
            "mealid": mealID,
        };
        
        doAjax(params, function(isFavourite) {
            if (isFavourite) {
                $(favButton).text("Unfavourite");
            }
            else {
                $(favButton).text("Favourite");
            }
        });
    };
    
    MPAjax.deleteMeal = function(mealID, rowElement) {
        var params = {
            "command": "delete-meal",
            "mealid": mealID,
        };
        
        doAjax(params, function(response) {
            $(rowElement).remove();
        });
    };
    
    return MPAjax;
})();
