// The MPUtil object provides routines that are used multiple times in the
// page-specific JS code. It mostly contains functions for rendering result sets
// as returned by MPAjax into HTML.

var MPUtil = (function() {
    var MPUtil = {};
    
    // Renders a single meal/tag result and returns the created <tr> element.
    // Used by MPUtil.renderMealList.
    function renderMealResult(result) {
        var row = $("<tr>");
        
        $("<td>").appendTo(row).text(result.meal.name);
        $("<td>").appendTo(row).text((result.tags || []).join(", "));
        var actions = $("<td>").addClass("action-buttons").appendTo(row);
        
        $("<button>")
            .appendTo(actions)
            .text("Open recipe")
            .click(function(event) {
                event.preventDefault();
                location.href = result.meal.recipe;
            });
        
        var favText = result.meal.favourite ? "Unfavourite" : "Favourite";
        $("<button>")
            .appendTo(actions)
            .text(favText)
            .click(function(event) {
                event.preventDefault();
                var favButton = $(this);
                
                MPAjax.toggleFavourite(result.meal.id, function(isFavourite) {
                    if (isFavourite) {
                        favButton.text("Unfavourite");
                    }
                    else {
                        favButton.text("Favourite");
                    }
                });
            });
        
        $("<button>")
            .appendTo(actions)
            .text("Edit")
            .click(function(event) {
                event.preventDefault();
                location.href = "/meals/" + result.meal.id + "/edit";
            });
        
        $("<button>")
            .appendTo(actions)
            .text("Delete")
            .click(function(event) {
                event.preventDefault();
                var row = $(this);
                
                MPAjax.deleteMeal(result.meal.id, function(response) {
                    row.remove();
                });
            });
        
        return row;
    }
    
    // Takes a list of meal/tag results (as returned by MPAjax.fetchMealList)
    // and renders them to a table created inside 'container'.
    MPUtil.renderMealList = function(results, container) {
        results = results || [];
        
        if (results.length == 0) {
            container.text("No meals to display.");
            return;
        }
        
        var table = $("<table>").appendTo(container.empty());
        var thead = $("<thead>").appendTo(table);
        var headerRow = $("<tr>").appendTo(thead);
        $("<th>").text("Name").appendTo(headerRow);
        $("<th>").text("Tags").appendTo(headerRow);
        $("<th>").text("Actions").appendTo(headerRow);
        var tbody = $("<tbody>").appendTo(table);
        
        var i, result, row;
        
        for (i = 0; i < results.length; i++) {
            row = renderMealResult(results[i]);
            row.appendTo(tbody);
        }
    };
    
    // Takes a list of tags (as returned by MPAjax.fetchAllTags) and renders
    // them to the <select> tag 'container'.
    MPUtil.renderExistingTagsList = function(tags, container) {
        tags = tags || [];
        
        var i, tag;
        for (i = 0; i < tags.length; i++) {
            tag = tags[i];
            $("<option>").val(tag).text(tag).appendTo(container);
        }
    }
    
    return MPUtil;
})();
