(function () {
  'use strict';

  angular.module('checklistApp', [ 'ngMaterial', 'ngMessages' ]);

  angular.module('checklistApp')
    .config(function ($mdThemingProvider, $mdIconProvider) {

      $mdIconProvider
        .icon('save', './bower_components/material-design-icons/content/svg/production/ic_save_24px.svg', 24)
      ;

      $mdThemingProvider.theme('default')
        .primaryPalette('blue')
        .accentPalette('pink');

    })

})();