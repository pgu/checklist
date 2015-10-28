(function () {
  'use strict';

  angular.module('checklistApp')
    .controller('ItemController', ItemController);

  function ItemController ($scope, $http, $timeout, $mdToast, $q) {

    var self = this;
    self.getNewItem = getNewItem;

    $scope.createItem = createItem;
    $scope.fetchItems = fetchItems;
    $scope.onChangeItem = onChangeItem;

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

          var items = response.data;
          $scope.items = _.sortBy(items, 'title');
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

          $scope.newItem = getNewItem();
          form.$setPristine();
          form.$setUntouched();

        })
        .then(fetchItems)

        .finally(function () {
          $scope.isSaving = false;
        });
    }

    function onChangeItem (item) {
      return $http.put('/items/' + item.id, item)
        .then(function () {
          $mdToast.showSimple('\u2713  ' + item.title);
        })
        .catch(function () {
          $mdToast.showSimple('\u274C  Could not be updated: ' + item.title);
          item.isDone != item.isDone;
        });
    }

  }

})();