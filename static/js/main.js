let btnId;
$(function () {
    $.ajax({
        url: '/api/v1/user/functions'
    }).done((r) => {
        $('#dynamicMenu').html(createMenu(r));

        $('.lang').on('click', function(e){
            e.preventDefault();
            const lang = $(this).data('lang');
            $.ajax({
                url: `/api/v1/lang?lang=${lang}`
            }).done(() => {
                location.reload();
            })
        });
    });
    if (btnId) {
        $('#' + btnId).css({'background-color': '#004783'});
    }

    $('.emptyTd').attr('colspan', $('th').length);

    $('a.page-link').on('click', function(e){
        e.preventDefault();
        const page = $(this).data('page');
        let search = location.search.slice(1);
        if(search.indexOf('page') === -1){
            if(search.length === 0){
                search = 'page=' + page;
            }else{
                search += '&page=' + page;
            }
        }else{
            search = search.replace(/page=\d+/, 'page=' + page)
        }
        location.assign(location.pathname + '?' + search);
    });

    $('#btnFilterSearch').on('click', function(e){
        e.preventDefault();
        const filterForm = $(this).closest('form');
        const queryStr = filterForm.find('input').filter((i, v) => v.value).serialize();
        const startDate = $('#searchStartDate').val();
        const endDate = $('#searchEndDate').val();
        if(startDate && endDate){
            if(startDate > endDate){
                BootstrapDialog.alert($('#strTimeError').text());
                return
            }
        }
        location.assign(location.pathname + (queryStr?'?':'') + queryStr);
    });

    $('#btnFilterReset').on('click', function(e){
        e.preventDefault();
        $(this).closest('form').find('input').val('');
    })
});

function createMenu(res) {
    let fullMenu = '<div>';
    let isLogout = false;
    if (res && res.roleId) {
        isLogout = true
    }

    fullMenu += createRightMenu(isLogout);
    if (!res || !res.functions || res.functions.length === 0 || filterFuncByParent(res.functions, '0').length === 0) {
        return fullMenu;
    }
    const menuIndex = getIdByFuncName(res.functions, 'menu');
    const menuItems = menuIndex ? filterFuncByParent(res.functions, menuIndex) : [];

    let block = '<div class="collapse navbar-collapse" id="navbar-collapse-1">' +
        '<ul class="nav navbar-nav">';
    let menu = menuItems.reduce((pre, cur) => {
        let html = '';
        let items = cur.id ? filterFuncByParent(res.functions, cur.id) : [];
        if (items.length === 0) {
            html = `<li><a href="${cur.uri}">${cur.name}</a></li>`;
        } else {
            html = '<li class="dropdown">' +
                `<a class="dropdown-toggle" role="button" data-toggle="dropdown" href="#" aria-haspopup="true" aria-expanded="false">${cur.name} ` +
                '<span class="caret"></span></a>' +
                '<ul class="dropdown-menu">' +
                items.reduce((p, c) => {
                    return p + `<li><a href="${c.uri}">${c.name}</a></li>`;
                }, '') +
                '</ul>' +
                '</li>';
        }
        return pre + html;
    }, '');

    fullMenu += block + menu;

    return fullMenu + '</div>';
}

function createRightMenu(isLogout){
    return '<ul class="nav navbar-nav navbar-right">' +
        '<li><a data-lang="zh_TW" class="lang" href="#">中文</a></li>' +
        '<li><a data-lang="zh_US" class="lang" href="#">English</a></li>' +
        (isLogout? '<li><a href="/logout">Logout</a></li>': '') +
        '</ul>';
}

function filterFuncByParent(functions, parent) {
    if (!functions) {
        return [];
    }
    return functions.filter((func) => {
        return func.parent == parent;
    })
}

function getIdByFuncName(functions, name) {
    if (!functions) {
        return '';
    }
    const func = functions.filter((func) => {
        return func.name == name;
    });
    return func.length === 0 ? '' : func[0].id;
}