// The MPAjax object provides functions for interacting with the server over
// Ajax.

var MPAjax = (function() {
    var MPAjax = {};
    
    // Generic function for performing an Ajax call with the given request data
    // and success callback.
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
    
    // Fetch the list of meals and add them in the form of a table to
    // 'destElement'.
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
            var thead = $("<thead>").appendTo(table);
            var headerRow = $("<tr>").appendTo(thead);
            $("<th>").text("Name").appendTo(headerRow);
            $("<th>").text("Tags").appendTo(headerRow);
            $("<th>").text("Actions").appendTo(headerRow);
            var tbody = $("<tbody>").appendTo(table);
            
            var i, result, row, actions, favText;
            
            for (i = 0; i < results.length; i++) {
                result = results[i];
                row = $("<tr>").appendTo(tbody);
                
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
    
    // Toggle the "favourite" status of the meal identified by 'mealID', and
    // update the text of 'favButton' to reflect the new status.
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
    
    // Delete the meal identified by 'mealID', and delete 'rowElement' if
    // successful.
    MPAjax.deleteMeal = function(mealID, rowElement) {
        var params = {
            "command": "delete-meal",
            "mealid": mealID,
        };
        
        doAjax(params, function(response) {
            $(rowElement).remove();
        });
    };
    
    MPAjax.tagsFetched = false;
    
    // Fetch the list of all tags in the database, and add them as <option>
    // elements to 'destElement'.
    MPAjax.fetchAllTags = function(destElement) {
        destElement = $(destElement);
        
        var params = {
            "command": "fetch-all-tags",
        };
        
        doAjax(params, function(tags) {
            tags = tags || [];
            
            var i, tag;
            for (i = 0; i < tags.length; i++) {
                tag = tags[i];
                $("<option>").val(tag).text(tag).appendTo(destElement);
            }
            
            MPAjax.tagsFetched = true;
        });
    };
    
    return MPAjax;
})();
