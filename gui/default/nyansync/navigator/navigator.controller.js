(function () {
  'use strict';

  angular
    .module('app')
    .controller('NavigatorController', ['$scope', 'NavigatorService',
      function ($scope, NavigatorService) {
        var vm = this;
        $scope.vm = vm;

        (function initController() {
        })();

        vm.navigator_path = [];
        vm.navigator_active = null;
        vm.navigator_active_el = null;

        vm.navigatePath = function(path) {
          vm.navigator_path = path;
          vm.navigator_items = NavigatorService.query(vm.navigator_path);
        };
        vm.navigateChildren = function(name) {
          vm.navigatePath(vm.navigator_path.concat([name]));
        };
        vm.itemMenu = function(item, e) {
          if( vm.navigator_active && vm.navigator_active !== item ) {
            $(vm.navigator_active_el).parent().children('.dropdown').remove();
          }
          vm.navigator_active = item;
          vm.navigator_active_el = e.currentTarget;
          var p = $(e.currentTarget).parent();
          if( p.children('.dropdown').length > 0 )
            return;

          // Create dropdown menu
          // TODO: replace with actual actions
          var el = $(
            '<div class="dropdown" style="position:absolute">'+
              '<button class="btn btn-default dropdown-toggle" type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">'+
                '<span class="caret"></span>'+
              '</button>'+
              '<ul class="dropdown-menu">'+
                '<li><a>Action</a></li>'+
                '<li><a>Another action</a></li>'+
                '<li><a>Something else here</a></li>'+
                '<li role="separator" class="divider"></li>'+
                '<li><a>Separated link</a></li>'+
              '</ul>'+
            '</div>');

          p.prepend(el);
        };
        vm.itemClick = function(item, e) {
          if( item === vm.navigator_active ) {
            if( item.Type == 'folder' )
              vm.navigateChildren(item.Name);
          } else {
            // For touch devices to use double click system
            if( vm.navigator_active_el ) {
              $(vm.navigator_active_el)
                .trigger('mouseout')
                .removeClass('active');
            }
            $(e.currentTarget)
              .trigger('mouseover')
              .addClass('active');

            vm.itemMenu(item, e);
          }
        };

        vm.navigatePath([]);
      }
    ]);

})();
