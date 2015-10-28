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
    $scope.deleteDoneItems = deleteDoneItems;

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


    function createItem (newItem, form, items) {

      if (form.$invalid) {
        form.title.$dirty = true;
        return;
      }

      $scope.isSaving = true;

      return $http.post('/items/', newItem)

        .then(function (response) {

          var strItemId = _.last(response.headers().location.split('/'));
          var itemId = _.parseInt(strItemId);

          var item = _.cloneDeep(newItem);
          item.id = itemId;
          items.push(item);

          $scope.items = _.sortBy(items, 'title');

          $scope.newItem = getNewItem();
          form.$setPristine();
          form.$setUntouched();

        })

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

    function deleteDoneItems (items) {

      var deletePromises = _(items)
        .filter({ isDone: true })
        .map(function (item) {
          return $http.delete('/items/' + item.id);
        })
        .value();

      $q.all(deletePromises)
        .then(function () {
          $scope.items = _.filter(items, { isDone: false });
        })
        .catch(function () {
          return fetchItems();
        });

    }

  }

})();