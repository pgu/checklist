(function () {
  'use strict';

  angular.module('checklistApp')
    .controller('ItemController', ItemController);

  function ItemController ($scope, $http, $timeout) {

    var self = this;
    self.getNewItem = getNewItem;

    $scope.createItem = createItem;
    $scope.fetchItems = fetchItems;

    // init
    $scope.newItem = self.getNewItem();
    fetchItems();


    function getNewItem () {
      return {
        title: '',
        isDone: false
      };
    }

    function fetchItems () {

      return $http.get('/items/')
        .then(function (response) {
          console.log(response);
          $scope.items = response.data;
        });

    }


    function createItem (newItem, form) {

      if (form.$invalid) {
        form.title.$dirty = true;
        return;
      }

      $scope.isSaving = true;

      return $http.post('/items/', newItem)

        .then(function () {

          $timeout(function () {

            $scope.newItem = getNewItem();
            form.$setPristine();

          }, 300);

          return fetchItems();
        })

        .finally(function () {
          $scope.isSaving = false;
        });
    }

  }

})();