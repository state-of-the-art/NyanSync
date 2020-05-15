(function() {
  'use strict';

  angular
    .module('app')
    .controller('AccessController', ['title', 'source_id', 'path', 'item', '$scope', '$uibModalInstance', '$uibModal', 'AccessService', 'UserService', 'SourceService',
      function( title, source_id, path, item, $scope, $uibModalInstance, $uibModal, AccessService, UserService, SourceService ) {
        var vm = this;

        vm.title = title;

        (function initController() {
          // Access to create or edit
          $scope.access = new AccessService();
          $scope.access.SourceId = source_id;
          $scope.access.Path = path;

          $scope.access_users = [];
          if( $scope.access.Users !== null ) {
            // TODO: fill access_users with actual objects
          }

          // To validate the source ids
          SourceService.query().$promise.then(function(sources){
            $scope.source_ids = sources.map(s => s.Id);
          });
        })();

        vm.submit = function() {
          $scope.access.$save().then(function(){
            $uibModalInstance.close();
          });
        };

        vm.cancel = function() {
          $uibModalInstance.dismiss('cancel');
        };

        vm.remove = function() {
          $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: 'nyanshare/modal/modal.question.html',
            controller: 'QuestionController',
            controllerAs: 'vm',
            size: 'sm',
            resolve: {
              body: function(){
                return 'Are you sure you want to remove access "' + $scope.access.Id + '"?';
              },
            },
          }).result.then(function( result ) {
            if( result === true ) {
              $scope.access.$remove().then(function(){
                $uibModalInstance.close();
              });
            }
          });
        };

        vm.loadUsers = function(query) {
          console.log(query);
          return UserService.query().$promise;
        };
        vm.createUser = function(data) {
          var promise = $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: 'nyanshare/modal/modal.user.html',
            controller: 'UserController',
            controllerAs: 'vm',
            size: 'md',
            resolve: {
              user: function(){ return data; },
            },
          }).result.then(function( result ) {
            console.log('asd', result);
            UserService.query({cache: false});
            if( result instanceof UserService ) {
              data.Login = result.Login
              data.Name = result.Name
            }
          });
          return promise;
        };
      }
    ])
    .directive('mySourceIdExist', function() {
      return {
        require: 'ngModel',
        scope: true,
        link: function(scope, elm, attrs, ctrl) {
          if( !ctrl )
            return;
          ctrl.$validators.sourceid = function(modelValue, viewValue) {
            return scope.source_ids.indexOf(modelValue) !== -1;
          };
        },
      };
    });

})();
