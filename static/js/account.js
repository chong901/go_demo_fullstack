let roleMap = {};
const ACTION_ADD = 'add';
const ACTION_EDIT = 'edit';
$(() => {
    $.ajax({
        url: '/api/v1/role/all'
    }).done((r) => {
        const data = r.data;
        if (data && data.length > 0) {
            data.forEach((role) => {
                roleMap[role.id] = role.name;
            });
        }
    }).fail(alertErr);
});

function openDialog(dom) {
    const $tr = $(dom).closest('tr');
    const roleId = $tr.find('.roleId').val();
    const account = $tr.find('.account').html();
    BootstrapDialog.show({
        title: `${$('.btnEditRole').val()}: ${account}`,
        size: BootstrapDialog.SIZE_SMALL,
        message: createDialogMessage(roleMap, roleId),
        buttons: [
            {
                label: $('#strBtnConfirm').text(),
                action: function (bd) {
                    const data = {
                        account,
                        roleId: parseInt($('#role').val())
                    };
                    $.ajax({
                        url: '/account/role',
                        method: 'PUT',
                        dataType: 'json',
                        data: JSON.stringify(data)
                    }).done((r) => {
                        const data = r.data;
                        BootstrapDialog.alert($('#strSaved').text());
                        setHtml($tr, 'roleName', roleMap[data.roleId]);
                        $tr.find('.roleId').val(data.roleId);
                    }).fail(alertErr).always(() => bd.close());
                }
            },
            {
                label: $('#strBtnCancel').text(),
                action: function (bd) {
                    bd.close();
                }
            }
        ]
    })
}

function createDialogMessage(roleMap, roleId) {
    let html = '<form>' +
        `<label for="role">${$('.roleStr').text()}</label>` +
        createRoleSelect(roleMap, roleId) +
        '</form>';
    return html;
}

function createRoleSelect(roleMap, val) {
    if (!roleMap) {
        return '<select id="role" class="form-control"></select>';
    }
    let html = '<select id="role" class="form-control">';
    for (let k in roleMap) {
        html += `<option value="${k}" ${k == val ? 'selected' : ''}>${roleMap[k]}</option>`;
    }
    html += '</select>';
    return html;
}

function userDialog(dom, action) {
    const $tr = $(dom).closest('tr');
    const account = $tr.find('.account').html();
    const infoObj = userDialogInfo(action);
    if (!infoObj) {
        BootstrapDialog.alert('Action is empty.');
        return;
    }
    BootstrapDialog.show({
        title: infoObj.title,
        message: createUserForm(account),
        onshown: function (bd) {
            $("#userForm").on('submit', function (ev) {
                ev.preventDefault();

                const data = getObjFromSerialize($(this).serializeArray());
                $.ajax({
                    url: infoObj.url,
                    method: infoObj.method,
                    dataType: 'json',
                    data: JSON.stringify(data)
                }).done((r) => {
                    BootstrapDialog.alert(infoObj.doneMsg);
                    if (action === ACTION_ADD) {
                        $('#dataTable').append(createUserTr(r.data || {}));
                        removeEmptyDataTr();
                    }
                    bd.close();
                }).fail(alertErr);
            });
        },
        buttons: [
            {
                label: $('#strBtnConfirm').text(),
                action: function () {
                    $("#userFormSubmit").click();
                }
            },
            {
                label: $('#strBtnCancel').text(),
                action: function (bd) {
                    bd.close();
                }
            }
        ]
    })
}

function createUserTr(result) {
    let tr = $('.hid tr').clone();
    tr.find('.account').text(result.account);
    tr.find('.roleName').text(roleMap[result.roleId]);
    tr.find('.roleId').val(result.roleId);
    return tr;
}

function userDialogInfo(action) {
    if (!action) {
        return;
    }
    let obj = {};
    switch (action) {
        case ACTION_ADD:
            obj.title = $('.btnAddAccount').val();
            obj.method = 'POST';
            obj.doneMsg = $('#strAccountAdded').text();
            obj.url = '/account';
            break;
        case ACTION_EDIT:
            obj.title = $('.btnEditPassword').val();
            obj.method = 'PUT';
            obj.doneMsg = $('#strPwdChanged').text();
            obj.url = '/account/password';
            break;
    }
    return obj;
}

function createUserForm(account) {
    let html = '<form id="userForm">' +
        '<div class="form-group"> ' +
        `<label for="account">${$('#strAccount').text()}</label> ` +
        `<input id="account" class="form-control" name="account" required value="${account || ''}" ${account ? 'readonly' : ''}> ` +
        '</div> ' +
        '<div class="form-group"> ' +
        `<label for="password">${$('#strPassword').text()}</label> ` +
        '<input id="password" class="form-control" type="password" required name="password"> ' +
        '</div> ' +
        '<input id="userFormSubmit" type="submit" style="display: none"> ' +
        '</form> ';
    return html;
}