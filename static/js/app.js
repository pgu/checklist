(function () {
  'use strict';

  angular.module('checklistApp', [ 'ngMaterial', 'ngMessages' ]);

  angular.module('checklistApp')
    .config(function ($mdThemingProvider) {

      $mdThemingProvider.theme('default')
        .primaryPalette('blue')
        .accentPalette('pink');

    })

})();