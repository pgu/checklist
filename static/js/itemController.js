(function () {
  'use strict';

  angular.module('checklistApp')
    .controller('ItemController', ItemController);

  function ItemController($scope, $http) {

    $scope.newItem = {
      isDone: false
    };

    $scope.createItem = createItem;
    $scope.fetchItems = fetchItems;

    function fetchItems() {

      return $http.get('/items/')
        .then(function (response) {
          console.log(response);
          $scope.items = response.data;
        });

    }


    function createItem(newItem) {
      return $http.post('/items/', newItem)
        .then(function () {
          $scope.newItem.title = '';
          return fetchItems();
        });
    }

  }

})();