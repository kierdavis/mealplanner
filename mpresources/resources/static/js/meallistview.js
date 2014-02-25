var MealResult = (function() {
    var MealResult = function(mt, score) {
        this.id = mt.meal.id;
        this.name = mt.meal.name;
        this.recipe = mt.meal.recipe;
        this.favourite = mt.meal.favourite;
        this.tags = mt.tags;
        this.score = score;
    };
    
    MealResult.fetchMealList = function(query, callback) {
        MPAjax.fetchMealList(query, function(mts) {
            var i, results = [];
            for (i = 0; i < mts.length; i++) {
                results.push(new MealResult(mts[i], null));
            }
            
            callback(results);
        });
    };
    
    MealResult.fetchSuggestions = function(mpID, date, callback) {
        MPAjax.fetchSuggestions(mpID, date, function(suggs) {
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
                view.itemCallback.call(view, item);
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
        var scoreStr = "" + MPUtil.round1dp(item.score * 9 + 1);
        if (scoreStr.indexOf(".") < 0) {
            scoreStr += ".0";
        }
        
        $("<td></td>").text(scoreStr).addClass(this.className).appendTo(row);
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
        var view = this.view;
        button.appendTo(cell).click(function(event) {
            event.preventDefault();
            
            if (confirm("Are you sure that you want to delete the meal '" + item.name + "'?")) {
                MPAjax.deleteMeal(item.id, function(response) {
                    view.deleteItemByID(item.id);
                });
            }
        });
    };
    
    return o;
})();

var MealListView = (function() {
    var MealListView = function(parent) {
        this.parent = parent;
        this.items = [];
        this.numPages = 0;
        this.currentPage = 0;
        this.columns = [];
        this.itemCallback = null;
        this.searchCallback = null;
        this.deleteCallback = null;
        this.tbody = null;
        this.pageNumSpan = null;
        this.numPagesSpan = null;
    };
    
    MealListView.prototype.setData = function(items) {
        this.items = items;
        this.currentPage = 0;
        this.touchData();
    };
    
    MealListView.prototype.touchData = function() {
        this.numPages = Math.floor((this.items.length + 9) / 10);
        this.renderCurrentPage();
    };
    
    MealListView.prototype.getCurrentPage = function() {
        return this.currentPage;
    };
    
    MealListView.prototype.setCurrentPage = function(p) {
        this.currentPage = p;
        this.renderCurrentPage();
    };
    
    MealListView.prototype.incrCurrentPage = function(amt) {
        this.currentPage += amt;
        this.renderCurrentPage();
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
        if (MPUtil.nonNull(this.deleteCallback)) {
            this.deleteCallback.call(this, this.items[idx]);
        }
        
        this.items.splice(idx, 1);
        this.touchData();
    };
    
    MealListView.prototype.highlightItemByID = function(id) {
        var idx = this.lookup(id);
        if (MPUtil.nonNull(idx)) {
            this.highlightItemByIndex(idx);
        }
    };
    
    MealListView.prototype.highlightItemByIndex = function(idx) {
        var newCurrentPage = Math.floor(idx / 10);
        if (this.currentPage != newCurrentPage) {
            this.setCurrentPage(newCurrentPage);
        }
        
        var row = $(this.tbody.find("tr")[idx % 10]);
        var bg = row.css("background");
        row.css("background", "orange");
        row.animate({backgroundColor: bg}, 1000);
    };
    
    MealListView.prototype.setItemCallback = function(cb) {
        this.itemCallback = cb;
    };
    
    MealListView.prototype.setSearchCallback = function(cb) {
        this.searchCallback = cb;
    };
    
    MealListView.prototype.setDeleteCallback = function(cb) {
        this.deleteCallback = cb;
    };
    
    MealListView.prototype.addColumn = function(col) {
        col.view = this;
        this.columns.push(col);
    };
    
    MealListView.prototype.render = function() {
        this.parent.empty();
        
        /*
        if (this.items.length == 0) {
            this.parent.text("No results to display.");
            return;
        }
        */
        
        this.renderSearch(this.parent);
        this.renderNav(this.parent);
        
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
        if (this.currentPage >= this.numPages) {
            this.currentPage = this.numPages - 1;
        }
        
        if (this.currentPage < 0) {
            this.currentPage = 0;
        }
        
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
    };
    
    MealListView.prototype.renderItem = function(item) {
        var row = $("<tr></tr>").appendTo(this.tbody);
        var i;
        for (i = 0; i < this.columns.length; i++) {
            this.columns[i].renderData(row, item);
        }
    };
    
    MealListView.prototype.renderSearch = function(parent) {
        if (MPUtil.nonNull(this.searchCallback)) {
            var container = $("<div class='table-search-container'></div>").appendTo(parent);
            var img = $("<img src='/static/img/loading.gif' height='16' alt='Searching...' style='margin-right: 10px'/>").appendTo(container).hide();
            var input = $("<input type='text' placeholder='Type to search...' width='30'/>").appendTo(container);
            
            var tid = null;
            var view = this;
            
            input.keydown(function() {
                if (MPUtil.nonNull(tid)) {
                    window.clearTimeout(tid);
                }
                
                tid = window.setTimeout(function() {
                    tid = null;
                
                    img.show();
                    view.searchCallback.call(view, input.val(), function() {
                        img.hide();
                    });
                }, 1200);
            });
        }
    };
    
    MealListView.prototype.renderNav = function(parent) {
        var nav = $("<div class='table-nav-container'></div>").appendTo(parent);
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
