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
                    alert("MPAjax error! Check console for more details.");
                }
                
                else {
                    if (MPUtil.nonNull(success)) {
                        success(resp.success);
                    }
                }
            },
        });
    }
    
    // Fetch the list of meals and pass them to a callback function.
    MPAjax.fetchMealList = function(callback) {
        var params = {
            "command": "fetch-meal-list",
        };
        
        doAjax(params, callback);
    };
    
    // Toggle the "favourite" status of the meal identified by 'mealID' and pass
    // the updated "favourite" status to a callback function.
    MPAjax.toggleFavourite = function(mealID, callback) {
        var params = {
            "command": "toggle-favourite",
            "mealid": mealID,
        };
        
        doAjax(params, callback);
    };
    
    // Delete the meal identified by 'mealID' and call the callback function
    // when done.
    MPAjax.deleteMeal = function(mealID, callback) {
        var params = {
            "command": "delete-meal",
            "mealid": mealID,
        };
        
        doAjax(params, callback);
    };
    
    // Fetch the list of all tags in the database and pass them to a callback
    // function.
    MPAjax.fetchAllTags = function(callback) {
        var params = {
            "command": "fetch-all-tags",
        };
        
        doAjax(params, callback);
    };
    
    // The first time this function is called, it is identical to 'fetchAllTags'.
    // After that, it does nothing.
    MPAjax.tagsFetched = false;
    MPAjax.fetchAllTagsOnce = function(callback) {
        if (!MPAjax.tagsFetched) {
            MPAjax.fetchAllTags(callback);
            MPAjax.tagsFetched = true;
        }
    };
    
    // Fetch a list of servings for the meal plan identified by 'mpID' and pass
    // them to a callback function.
    MPAjax.fetchServings = function(mpID, callback) {
        var params = {
            "command": "fetch-servings",
            "mealplanid": mpID,
        };
        
        doAjax(params, callback);
    };
    
    MPAjax.fetchSuggestions = function(date, callback) {
        var params = {
            "command": "fetch-suggestions",
            "date": MPUtil.formatDateJSON(date),
        };
        
        doAjax(params, callback);
    }
    
    MPAjax.updateServing = function(mpID, date, mealID, callback) {
        var params = {
            "command": "update-serving",
            "mealplanid": mpID,
            "date": MPUtil.formatDateJSON(date),
            "mealid": mealID,
        };
        
        doAjax(params, callback);
    };
    
    MPAjax.deleteServing = function(mpID, date, callback) {
        var params = {
            "command": "delete-serving",
            "mealplanid": mpID,
            "date": MPUtil.formatDateJSON(date),
        };
        
        doAjax(params, callback);
    };
    
    MPAjax.updateNotes = function(mpID, notes, callback) {
        var params = {
            "command": "update-notes",
            "mealplanid": mpID,
            "notes": notes,
        };
        
        doAjax(params, callback);
    }
    
    return MPAjax;
})();
