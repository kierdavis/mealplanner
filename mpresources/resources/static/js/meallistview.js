var MealResult = (function() {
    var MealResult = function(mt, score) {
        this.id = mt.meal.id;
        this.name = mt.meal.name;
        this.recipe = mt.meal.recipe;
        this.favourite = mt.meal.favourite;
        this.tags = mt.tags;
        this.score = score;
    };
    
    MealResult.fetchMealList = function(callback) {
        MPAjax.fetchMealList(function(mts) {
            var i, results = [];
            for (i = 0; i < mts.length; i++) {
                results.push(new MealResult(mts[i], null));
            }
            
            callback(results);
        });
    };
    
    MealResult.fetchSuggestions = function(date, callback) {
        MPAjax.fetchSuggestions(date, function(suggs) {
            var i, results = [];
            for (i = 0; i < suggs.length; i++) {
                results.push(new MealResult(suggs[i].mt, suggs[i].score));
            }
            
            callback(results);
        });
    };
    
    MealResult.prototype.hasScore = function() {
        return MPUtil.nonNull(this.score);
    };
    
    return MealResult;
})();

var MealListViewColumns = (function() {
    var o = {};
    
    o.NameColumn = function(className) {
        this.view = null;
        this.className = className;
    };
    o.NameColumn.prototype.renderHeader = function(row) {
        $("<th>Name</th>").addClass(this.className).appendTo(row);
    };
    o.NameColumn.prototype.renderData = function(row, item) {
        var cell = $("<td></td>").addClass(this.className).appendTo(row);
        
        if (MPUtil.nonNull(this.view.itemCallback)) {
            var view = this.view;
            var link = $("<a href='#'></a>").text(item.name).appendTo(cell).click(function(event) {
                event.preventDefault();
                view.itemCallback(item);
            });
        }
        
        else {
            cell.text(item.name);
        }
    };
    
    o.TagsColumn = function(className) {
        this.view = null;
        this.className = className;
    };
    o.TagsColumn.prototype.renderHeader = function(row) {
        $("<th>Tags</th>").addClass(this.className).appendTo(row);
    };
    o.TagsColumn.prototype.renderData = function(row, item) {
        var tagsString = (item.tags || []).join(", ");
        $("<td></td>").text(tagsString).addClass(this.className).appendTo(row);
    };
    
    o.ScoreColumn = function(className) {
        this.view = null;
        this.className = className;
    };
    o.ScoreColumn.prototype.renderHeader = function(row) {
        $("<th>Score</th>").addClass(this.className).appendTo(row);
    };
    o.ScoreColumn.prototype.renderData = function(row, item) {
        $("<td></td>").text(item.score).addClass(this.className).appendTo(row);
    };
    
    o.ActionsColumn = function(individualClassName, spannedClassName) {
        this.view = null;
        this.individualClassName = individualClassName;
        this.spannedClassName = spannedClassName;
    };
    o.ActionsColumn.prototype.renderHeader = function(row) {
        $("<th colspan='4'>Actions</th>").addClass(this.spannedClassName).appendTo(row);
    };
    o.ActionsColumn.prototype.renderData = function(row, item) {
        this.renderRecipeButton(row, item);
        this.renderFavButton(row, item);
        this.renderEditButton(row, item);
        this.renderDeleteButton(row, item);
    };
    
    o.ActionsColumn.prototype.renderRecipeButton = function(row, item) {
        var cell = $("<td></td>").addClass(this.individualClassName).appendTo(row);
        
        if (item.recipe) {
            var button = $("<button title='Open the recipe page listed for this meal' class='action-button'><img src='/static/img/open-recipe_16x16.png' height='16' alt=''/></button>");
            button.appendTo(cell).click(function(event) {
                event.preventDefault();
                location.href = item.recipe;
            });
        }
    };
    
    o.ActionsColumn.prototype.renderFavButton = function(row, item) {
        var cell = $("<td></td>").addClass(this.individualClassName).appendTo(row);
        
        var favButton   = $("<button title='Mark this meal as a favourite' class='action-button'><img src='/static/img/favourite_16x16.png' height='16' alt=''/></button>");
        var unfavButton = $("<button title='Remove the favourite mark from this meal' class='action-button'><img src='/static/img/unfavourite_16x16.png' height='16' alt=''/></button>");
        
        var toggleFavCallback = function(event) {
            event.preventDefault();
            MPAjax.toggleFavourite(item.id, function(isFavourite) {
                item.favourite = isFavourite;
                
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
        
        favButton.appendTo(cell).click(toggleFavCallback);
        unfavButton.appendTo(cell).click(toggleFavCallback);
        
        if (item.favourite) {
            favButton.hide();
        }
        else {
            unfavButton.hide();
        }
    };
    
    o.ActionsColumn.prototype.renderEditButton = function(row, item) {
        var cell = $("<td></td>").addClass(this.individualClassName).appendTo(row);
        var button = $("<button title='Edit this meal' class='action-button'><img src='/static/img/edit_24x24.png' height='16' alt=''/></button>");
        button.appendTo(cell).click(function(event) {
            event.preventDefault();
            location.href = "/meals/" + item.id + "/edit";
        });
    };
    
    o.ActionsColumn.prototype.renderDeleteButton = function(row, item) {
        var cell = $("<td></td>").addClass(this.individualClassName).appendTo(row);
        var button = $("<button title='Delete this meal from the database' class='action-button'><img src='/static/img/delete_24x24.png' height='16' alt=''/></button>");
        button.appendTo(cell).click(function(event) {
            event.preventDefault();
            
            if (confirm("Are you sure that you want to delete the meal '" + item.name + "'?")) {
                MPAjax.deleteMeal(item.id, function(response) {
                    this.view.deleteItemByID(item.id);
                });
            }
        });
    };
    
    return o;
})();

var MealListView = (function() {
    var MealListView = function(parent, items) {
        this.parent = parent;
        this.items = items;
        this.numPages = Math.floor((items.length + 9) / 10);
        this.currentPage = 0;
        this.columns = [];
        this.itemCallback = null;
        this.tbody = null;
        this.pageNumSpan = null;
        this.numPagesSpan = null;
        this.highlightRowNum = null;
    };
    
    MealListView.prototype.getCurrentPage = function() {
        return this.currentPage;
    };
    
    MealListView.prototype.setCurrentPage = function(p) {
        if (p < 0) {
            p = 0;
        }
        
        if (p >= this.numPages) {
            p = this.numPages - 1;
        }
        
        this.currentPage = p;
        this.renderCurrentPage();
    };
    
    MealListView.prototype.incrCurrentPage = function(amt) {
        this.setCurrentPage(this.currentPage + amt);
    };
    
    MealListView.prototype.lookup = function(id) {
        var i;
        for (i = 0; i < this.items.length; i++) {
            if (this.items[i].id == id) {
                return i;
            }
        }
        return null;
    };
    
    MealListView.prototype.deleteItemByID = function(id) {
        var idx = this.lookup(id);
        if (MPUtil.nonNull(idx)) {
            this.deleteItemByIndex(idx);
        }
    };
    
    MealListView.prototype.deleteItemByIndex = function(idx) {
        this.items.splice(idx, 1);
        
        if (this.items.length == 0) {
            this.parent.text("No results to display.");
            return;
        }
        
        // Update the number of pages.
        this.numPages = Math.floor((items.length + 9) / 10);
        
        // Check that currentPage is within the new bounds, and redraw the list.
        this.setCurrentPage(this.currentPage);
    };
    
    MealListView.prototype.highlightItemByID = function(id) {
        var idx = this.lookup(id);
        if (MPUtil.nonNull(idx)) {
            this.highlightItemByIndex(idx);
        }
    };
    
    MealListView.prototype.highlightItemByIndex = function(idx) {
        this.currentPage = Math.floor(idx / 10);
        this.highlightRowNum = idx % 10;
    };
    
    MealListView.prototype.setItemCallback = function(cb) {
        this.itemCallback = cb;
    };
    
    MealListView.prototype.addColumn = function(col) {
        col.view = this;
        this.columns.push(col);
    };
    
    MealListView.prototype.render = function() {
        this.parent.empty();
        
        if (this.items.length == 0) {
            this.parent.text("No results to display.");
            return;
        }
        
        this.renderNav($("<div style='width: 100%'></div>").appendTo(this.parent));
        
        var table = $("<table style='width: 100%; table-layout: fixed'></table>").appendTo(this.parent);
        var thead = $("<thead></thead>").appendTo(table);
        var tbody = $("<tbody></tbody>").appendTo(table);
        
        var headerRow = $("<tr></tr>").appendTo(thead);
        var i;
        for (i = 0; i < this.columns.length; i++) {
            this.columns[i].renderHeader(headerRow);
        }
        
        this.tbody = tbody;
        this.renderCurrentPage();
    };
    
    MealListView.prototype.renderCurrentPage = function() {
        this.tbody.empty();
        
        this.pageNumSpan.text(this.currentPage + 1);
        this.numPagesSpan.text(this.numPages);
        
        var start = this.currentPage * 10;
        var end = start + 10;
        if (end > this.items.length) {
            end = this.items.length;
        }
        
        var i;
        for (i = start; i < end; i++) {
            this.renderItem(this.items[i]);
        }
        
        if (MPUtil.nonNull(this.highlightRowNum)) {
            var row = $(tbody.find("tr")[this.highlightRowNum]);
            var bg = row.css("background");
            row.css("background", "orange");
            row.animate({
                backgroundColor: bg,
            }, 1000);
            
            this.highlightRowNum = null;
        }
    };
    
    MealListView.prototype.renderItem = function(item) {
        var row = $("<tr></tr>").appendTo(this.tbody);
        var i;
        for (i = 0; i < this.columns.length; i++) {
            this.columns[i].renderData(row, item);
        }
    };
    
    MealListView.prototype.renderNav = function(nav) {
        var left = $("<div class='table-nav table-nav-left'></div>").appendTo(nav);
        var center = $("<div class='table-nav table-nav-center'></div>").appendTo(nav);
        var right = $("<div class='table-nav table-nav-right'></div>").appendTo(nav);
        
        var view = this;
        
        var firstButton = $("<button title='Navigate to the first page of results'><img src='/static/img/first_24x24.png' height='16' alt='First'/></button>");
        firstButton.appendTo(left).click(function(event) {
            event.preventDefault();
            view.setCurrentPage(0);
        });
        
        var prevButton = $("<button title='Navigate to the previous page of results'><img src='/static/img/prev_24x24.png' height='16' alt='Prev'/></button>");
        prevButton.appendTo(left).click(function(event) {
            event.preventDefault();
            view.incrCurrentPage(-1);
        });
        
        var nextButton = $("<button title='Navigate to the next page of results'><img src='/static/img/next_24x24.png' height='16' alt='Next'/></button>");
        nextButton.appendTo(right).click(function(event) {
            event.preventDefault();
            view.incrCurrentPage(1);
        });
        
        var lastButton = $("<button title='Navigate to the last page of results'><img src='/static/img/last_24x24.png' height='16' alt='Last'/></button>");
        lastButton.appendTo(right).click(function(event) {
            event.preventDefault();
            view.setCurrentPage(view.numPages - 1);
        });
        
        center.html("Page <span id='page-num'></span> of <span id='num-pages'></span>");
        this.pageNumSpan = center.find("#page-num");
        this.numPagesSpan = center.find("#num-pages");
    };
    
    return MealListView;
})();
