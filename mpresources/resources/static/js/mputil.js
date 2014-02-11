// The MPUtil object provides routines that are used multiple times in the
// page-specific JS code. It mostly contains functions for rendering result sets
// as returned by MPAjax into HTML.

var MPUtil = (function() {
    var shortWeekdays = [
        "Sun",
        "Mon",
        "Tue",
        "Wed",
        "Thu",
        "Fri",
        "Sat",
    ];
    
    var shortMonths = [
        "Jan",
        "Feb",
        "Mar",
        "Apr",
        "May",
        "Jun",
        "Jul",
        "Aug",
        "Sep",
        "Oct",
        "Nov",
        "Dec",
    ];
    
    var MPUtil = {};
    
    function zeroPad(str, length) {
        str = "" + str;
        while (str.length < length) {
            str = "0" + str;
        }
        return str;
    }
    
    // Renders a single meal/tag result and returns the created <tr> element.
    // Used by MPUtil.renderMealList.
    function renderMealResult(mt, score, callback) {
        var row = $("<tr>");
        var nameCell = $("<td>").appendTo(row);
        
        if (MPUtil.nonNull(callback)) {
            $("<a href='#'>").text(mt.meal.name).appendTo(nameCell).click(function(event) {
                event.preventDefault();
                callback(mt);
            });
        }
        else {
            nameCell.text(mt.meal.name);
        }
        
        $("<td>").appendTo(row).text((mt.tags || []).join(", "));
        
        if (MPUtil.nonNull(score)) {
            $("<td>").text(score).appendTo(row);
        }
        
        if (mt.meal.recipe) {
            $("<td><button title='Open the recipe page listed for this meal' class='action-button'><img src='/static/img/open-recipe_16x16.png' height='16' alt=''/></button></td>")
                .appendTo(row)
                .click(function(event) {
                    event.preventDefault();
                    location.href = mt.meal.recipe;
                });
        }
        else {
            $("<td>").appendTo(row);
        }
        
        var favButton   = $("<button title='Mark this meal as a favourite' class='action-button'><img src='/static/img/favourite_16x16.png' height='16' alt=''/></button>");
        var unfavButton = $("<button title='Remove the favourite mark from this meal' class='action-button'><img src='/static/img/unfavourite_16x16.png' height='16' alt=''/></button>");
        
        favButton.click(function(event) {
            event.preventDefault();
            
            MPAjax.toggleFavourite(mt.meal.id, function(isFavourite) {
                if (isFavourite) {
                    favButton.hide();
                    unfavButton.show();
                }
                else {
                    alert("Unexpected: 'favourite' button did not favourite the meal!");
                }
            });
        });
        
        unfavButton.click(function(event) {
            event.preventDefault();
            
            MPAjax.toggleFavourite(mt.meal.id, function(isFavourite) {
                if (!isFavourite) {
                    unfavButton.hide();
                    favButton.show();
                }
                else {
                    alert("Unexpected: 'unfavourite' button did not unfavourite the meal!");
                }
            })
        });
        
        if (mt.meal.favourite) {
            favButton.hide();
        }
        else {
            unfavButton.hide();
        }
        
        $("<td>").appendTo(row).append(favButton).append(unfavButton);
        
        $("<td><button title='Edit this meal' class='action-button'><img src='/static/img/edit_24x24.png' height='16' alt=''/></button></td>")
            .appendTo(row)
            .click(function(event) {
                event.preventDefault();
                location.href = "/meals/" + mt.meal.id + "/edit";
            });
        
        $("<td><button title='Delete this meal from the list' class='action-button'><img src='/static/img/delete_24x24.png' height='16' alt=''/></button></td>")
            .appendTo(row)
            .click(function(event) {
                event.preventDefault();
                
                if (confirm("Are you sure you want to delete the meal '" + mt.meal.name + "'?")) {
                    MPAjax.deleteMeal(mt.meal.id, function(response) {
                        row.remove();
                    });
                }
            });
        
        return row;
    }
    
    // Takes a list of meal/tag results (as returned by MPAjax.fetchMealList)
    // and renders them to a table created inside 'container'. 'callback', if
    // not null, is a function that will be called when the meal name is clicked.
    // It is passed the meal-with-tags object.
    MPUtil.renderMealList = function(mts, container, callback) {
        mts = mts || [];
        container.empty();
        
        if (mts.length == 0) {
            container.text("No results to display.");
            return;
        }
        
        var table = $("<table>").appendTo(container);
        var thead = $("<thead>").appendTo(table);
        var headerRow = $("<tr>").appendTo(thead);
        $("<th>Name</th>").appendTo(headerRow);
        $("<th>Tags</th>").appendTo(headerRow);
        $("<th colspan='4'>Actions</th>").appendTo(headerRow);
        var tbody = $("<tbody>").appendTo(table);
        
        var i, result, row;
        
        for (i = 0; i < mts.length; i++) {
            row = renderMealResult(mts[i], null, callback);
            row.appendTo(tbody);
        }
    };
    
    MPUtil.renderSuggestions = function(suggs, container, callback) {
        suggs = suggs || [];
        container.empty();
        
        if (suggs.length == 0) {
            container.text("No results to display.");
            return;
        }
        
        var table = $("<table>").appendTo(container);
        var thead = $("<thead>").appendTo(table);
        var headerRow = $("<tr>").appendTo(thead);
        $("<th>Name</th>").appendTo(headerRow);
        $("<th>Tags</th>").appendTo(headerRow);
        $("<th>Score</th>").appendTo(headerRow);
        $("<th colspan='4'>Actions</th>").appendTo(headerRow);
        var tbody = $("<tbody>").appendTo(table);
        
        var i, result, row;
        
        for (i = 0; i < suggs.length; i++) {
            row = renderMealResult(suggs[i].mt, 1.0 * suggs[i].score, callback);
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
    };
    
    MPUtil.formatMonthHumanReadable = function(date) {
        return shortMonths[date.getMonth()] + " " + date.getFullYear();
    };
    
    MPUtil.formatDateHumanReadable = function(date) {
        return shortWeekdays[date.getDay()] + " " + date.getDate() + " " + shortMonths[date.getMonth()];
    };
    
    MPUtil.formatDateJSON = function(date) {
        return zeroPad(date.getFullYear(), 4) + "-" + zeroPad(date.getMonth() + 1, 2) + "-" + zeroPad(date.getDate(), 2);
    };
    
    MPUtil.nonNull = function(value) {
        return typeof value !== "undefined" && value !== null;
    };
    
    return MPUtil;
})();
