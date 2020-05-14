(function () {
  'use strict';

  angular.module('app')
    .factory('DropdownService', [function () {
      var base = [
        '<div class="dropdown" style="position:absolute">',
          '<button class="btn btn-default dropdown-toggle" type="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">',
            '<span class="caret"></span>',
          '</button>',
          '<ul class="dropdown-menu">',
          '</ul>',
        '</div>',
      ].join('');

      function create(actions, el) {
        var html = $(base);
        var ul = html.children('.dropdown-menu');

        for( var a in actions ) {
          var a_el = $('<a/>').click(actions[a]).text(a);
          var li = $('<li/>').append(a_el);
          ul.append(li);
        }

        el.prepend(html);
      };

      function remove(el) {
        el.children('.dropdown').remove();
      };

      function exist(el) {
        return el.children('.dropdown').length > 0;
      };

      return {
        Create: create,
        Remove: remove,
        Exist: exist,
      }
    }]);

})();
