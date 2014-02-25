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
    
    function renderNameCell(mt, callback) {
        var nameCell = $("<td>");
        
        if (MPUtil.nonNull(callback)) {
            $("<a href='#'>").text(mt.meal.name).appendTo(nameCell).click(function(event) {
                event.preventDefault();
                callback(mt);
            });
        }
        else {
            nameCell.text(mt.meal.name);
        }
        
        return nameCell;
    }
    
    function renderTagsCell(mt) {
        return $("<td>").text((mt.tags || []).join(", "));
    }
    
    function renderScoreCell(score) {
        return $("<td>").text(score);
    }
    
    function renderRecipeCell(mt) {
        if (mt.meal.recipe) {
            var button = $("<button title='Open the recipe page listed for this meal' class='action-button'><img src='/static/img/open-recipe_16x16.png' height='16' alt=''/></button>")
                .click(function(event) {
                    event.preventDefault();
                    location.href = mt.meal.recipe;
                });
            
            return $("<td>").append(button);
        }
        else {
            return $("<td>");
        }
    }
    
    function renderFavCell(mt) {
        var toggleFavCallback = function(event) {
            event.preventDefault();
            
            MPAjax.toggleFavourite(mt.meal.id, function(isFavourite) {
                if (isFavourite) {
                    favButton.hide();
                    unfavButton.show();
                }
                else {
                    unfavButton.hide();
                    favButton.show();
                }
            });
        };
        
        var favButton   = $("<button title='Mark this meal as a favourite' class='action-button'><img src='/static/img/favourite_16x16.png' height='16' alt=''/></button>");
        var unfavButton = $("<button title='Remove the favourite mark from this meal' class='action-button'><img src='/static/img/unfavourite_16x16.png' height='16' alt=''/></button>");
        
        favButton.click(toggleFavCallback);
        unfavButton.click(toggleFavCallback);
        
        if (mt.meal.favourite) {
            favButton.hide();
        }
        else {
            unfavButton.hide();
        }
        
        return $("<td>").append(favButton).append(unfavButton);
    }
    
    function renderEditCell(mt) {
        return $("<td><button title='Edit this meal' class='action-button'><img src='/static/img/edit_24x24.png' height='16' alt=''/></button></td>")
            .click(function(event) {
                event.preventDefault();
                location.href = "/meals/" + mt.meal.id + "/edit";
            });
    }
    
    function renderDeleteCell(mt, row) {
        return $("<td><button title='Delete this meal from the database' class='action-button'><img src='/static/img/delete_24x24.png' height='16' alt=''/></button></td>")
            .click(function(event) {
                event.preventDefault();
                
                if (confirm("Are you sure you want to delete the meal '" + mt.meal.name + "'?")) {
                    MPAjax.deleteMeal(mt.meal.id, function(response) {
                        location.reload();
                    });
                }
            });
    }
    
    // Renders a single meal/tag result and returns the created <tr> element.
    // Used by MPUtil.renderMealList.
    function renderMealRow(mt, callback) {
        var row = $("<tr>");
        
        row.append(renderNameCell(mt, callback).addClass("meal-list-name"));
        row.append(renderTagsCell(mt).addClass("meal-list-tags"));
        row.append(renderRecipeCell(mt).addClass("meal-list-action"));
        row.append(renderFavCell(mt).addClass("meal-list-action"));
        row.append(renderEditCell(mt).addClass("meal-list-action"));
        row.append(renderDeleteCell(mt, row).addClass("meal-list-action"));
        
        return row;
    }
    
    // Renders a single meal/tag result and returns the created <tr> element.
    // Used by MPUtil.renderMealList.
    function renderSuggRow(sugg, callback) {
        var row = $("<tr>");
        
        row.append(renderNameCell(sugg.mt, callback).addClass("sugg-list-name"));
        row.append(renderTagsCell(sugg.mt).addClass("sugg-list-tags"));
        row.append(renderScoreCell(sugg.score).addClass("sugg-list-score"));
        row.append(renderRecipeCell(sugg.mt).addClass("sugg-list-action"));
        row.append(renderFavCell(sugg.mt).addClass("sugg-list-action"));
        row.append(renderEditCell(sugg.mt).addClass("sugg-list-action"));
        row.append(renderDeleteCell(sugg.mt).addClass("sugg-list-action"));
        
        return row;
    }
    
    function renderPage(items, tbody, renderRow, callback, start, count) {
        tbody.empty();
        
        var end = start + count;
        if (end > items.length) end = items.length;
        
        var i, row;
        var alt = true;
        
        for (i = start; i < end; i++) {
            row = renderRow(items[i], callback);
            if (alt) row.addClass("alt");
            row.appendTo(tbody);
            
            alt = !alt;
        }
    }
    
    function updatePageNumCell(pageNumCell, page, numPages) {
        pageNumCell.text("Page " + (page + 1) + " of " + numPages);
    }
    
    function renderPaged(items, container, numCols, headerRow, renderRow, callback, highlightPred) {
        items = items || [];
        container.empty();
        
        if (items.length == 0) {
            container.text("No results to display.");
            return;
        }
        
        var numPages = Math.floor((items.length + 9) / 10);
        
        var table = $("<table class='meal-list'>").appendTo(container);
        var thead = $("<thead>").appendTo(table);
        var tbody = $("<tbody>").appendTo(table);
        
        var navCell = $("<td colspan='"+numCols+"' style='padding-bottom: 20px'>").appendTo($("<tr>").appendTo(thead));
        headerRow.appendTo(thead);
        
        var page = 0;
        var highlightIndex = null;
        
        if (MPUtil.nonNull(highlightPred)) {
            var i;
            for (i = 0; i < items.length; i++) {
                if (highlightPred(items[i])) {
                    page = Math.floor(i / 10);
                    highlightIndex = i % 10;
                    break;
                }
            }
        }
        
        var prevCell    = $("<div class='table-nav table-nav-left'>").appendTo(navCell);
        var pageNumCell = $("<div class='table-nav table-nav-center'>").appendTo(navCell);
        var nextCell    = $("<div class='table-nav table-nav-right'>").appendTo(navCell);
        
        $("<button title='Navigate to the first page of results'><img src='/static/img/first_24x24.png' height='16' alt='First'/></button>")
            .appendTo(prevCell)
            .click(function(event) {
                event.preventDefault();
                
                page = 0;
                updatePageNumCell(pageNumCell, page, numPages);
                renderPage(items, tbody, renderRow, callback, page * 10, 10);
            });
        
        $("<button title='Navigate to the previous page of results'><img src='/static/img/prev_24x24.png' height='16' alt='Prev'/></button>")
            .appendTo(prevCell)
            .click(function(event) {
                event.preventDefault();
                
                if ((page - 1) >= 0) {
                    page -= 1;
                }
                
                updatePageNumCell(pageNumCell, page, numPages);
                renderPage(items, tbody, renderRow, callback, page * 10, 10);
            });
        
        $("<button title='Navigate to the next page of results'><img src='/static/img/next_24x24.png' height='16' alt='Next'/></button>")
            .appendTo(nextCell)
            .click(function(event) {
                event.preventDefault();
                
                if ((page + 1) < numPages) {
                    page += 1;
                }
                
                updatePageNumCell(pageNumCell, page, numPages);
                renderPage(items, tbody, renderRow, callback, page * 10, 10);
            });
        
        $("<button title='Navigate to the last page of results'><img src='/static/img/last_24x24.png' height='16' alt='Last'/></button>")
            .appendTo(nextCell)
            .click(function(event) {
                event.preventDefault();
                
                page = numPages - 1;
                
                updatePageNumCell(pageNumCell, page, numPages);
                renderPage(items, tbody, renderRow, callback, page * 10, 10);
            });
        
        updatePageNumCell(pageNumCell, page, numPages);
        renderPage(items, tbody, renderRow, callback, page * 10, 10);
        
        if (MPUtil.nonNull(highlightIndex)) {
            var row = $(tbody.find("tr")[highlightIndex]);
            var bg = row.css("background");
            row.css("background", "#0f0");
            row.animate({
                backgroundColor: bg,
            }, 1000);
        }
    };
    
    // Takes a list of meal/tags results (as returned by MPAjax.fetchMealList)
    // and renders them to a table created inside 'container'. 'callback', if
    // not null, is a function that will be called when the meal name is clicked.
    // It is passed the meal/tags object.
    MPUtil.renderMealList = function(mts, container, callback) {
        var headerRow = $("<tr><th class='meal-list-name'>Name</th><th class='meal-list-tags'>Tags</th><th colspan='4' class='meal-list-actions'>Actions</th></tr>");
        renderPaged(mts, container, 6, headerRow, renderMealRow, callback, null);
    };
    
    MPUtil.renderMealListHighlight = function(mts, container, highlightID, callback) {
        var headerRow = $("<tr><th class='meal-list-name'>Name</th><th class='meal-list-tags'>Tags</th><th colspan='4' class='meal-list-actions'>Actions</th></tr>");
        var highlightPred = function(mt) {return mt.meal.id == highlightID};
        renderPaged(mts, container, 6, headerRow, renderMealRow, callback, highlightPred);
    };
    
    MPUtil.renderSuggestions = function(suggs, container, callback) {
        var headerRow = $("<tr><th class='sugg-list-name'>Name</th><th class='sugg-list-tags'>Tags</th><th class='sugg-list-score'>Score</th><th colspan='4' class='sugg-list-actions'>Actions</th></tr>")
        renderPaged(suggs, container, 7, headerRow, renderSuggRow, callback, null);
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
    
    MPUtil.parseDatepickerDate = function(str) {
        parts = str.split("/");
        if (parts.length < 3 || 1*parts[2] == NaN || 1*parts[1] == NaN || 1*parts[0] == NaN) {
            return null;
        }
        
        return new Date(parts[2], parts[1]-1, parts[0]);
    };
    
    MPUtil.nonNull = function(value) {
        return typeof value !== "undefined" && value !== null;
    };
    
    MPUtil.round1dp = function(x) {
        return Math.round(x * 10) / 10;
    };
    
    return MPUtil;
})();
