$(() => {
    $('#roleSelect').change(function () {
        const $this = $(this);

        if (!$this.val()) return;

        $.ajax({
            url: `/api/v1/role?id=${$this.val()}`
        }).done((r) => {
            const roleFuncs = r.data.functions || [];
            const roleFuncsIds = roleFuncs.map((r) => r.functionId);

            if ($('#functionData input:checkbox').length === 0) {
                return;
            }

            const checked = $('#functionData input:checkbox').filter(function () {
                return roleFuncsIds.indexOf(parseInt($(this).val())) !== -1;
            });
            const unchecked = $('#functionData input:checkbox').not(checked);
            checked.prop('checked', true);
            unchecked.prop('checked', false);
            $('#checkAll').prop('checked', false);
        }).fail(alertErr);
    });

    $.ajax({
        url: '/api/v1/role/all'
    }).done((r) => {
        const data = r.data || [];
        setRoleSelect(data);
        $('#roleSelect').trigger('change');
    }).fail(alertErr);

});

function setRoleSelect(data) {
    $('#roleSelect').html(data.reduce((pre, cur) => {
        var temp = `<option value="${cur.id || ''}">${cur.name || ''}</option>`;
        return pre + temp;
    }, ''));
}

function editRole(flag) {
    if (flag === 'add') {
        createRoleDialog();
    } else {
        const $roleSelect = $('#roleSelect');
        const id = $roleSelect.val();
        const name = $roleSelect.find(`[value="${id}"]`).html();
        createRoleDialog(id, name);
    }
}

function createRoleDialog(id, name) {
    BootstrapDialog.show({
        title: id ? $('#strBtnEdit').text() : $('#strAdd').text(),
        message: createRoleForm(id, name),
        buttons: [
            {
                label: $('#strBtnConfirm').text(),
                action: function (bd) {
                    const data = getObjFromSerialize($('#roleForm').serializeArray(), ['id']);
                    $.ajax({
                        url: id ? '/role/name' : '/role',
                        method: id ? 'PUT' : 'POST',
                        dataType: 'json',
                        data: JSON.stringify(data)
                    })
                        .done((r) => {
                            alterRoleSelect(r.data);
                            bd.close();
                        })
                        .fail(alertErr);
                }
            },
            {
                label: $('#strBtnCancel').text(),
                action: function (bd) {
                    bd.close();
                }
            }
        ]
    });
}

function alterRoleSelect(data) {
    if (!data) {
        return;
    }
    const $roleSelect = $('#roleSelect');
    const $option = $roleSelect.find(`[value="${data.id}"]`);
    if ($option.length == 0) {
        $roleSelect.append(`<option value="${data.id}">${data.name}</option>`);
    } else {
        $option.html(data.name);
    }
}

function createRoleForm(id, name) {
    let html = '<form id="roleForm">' +
        '<div class="form-group">' +
        `<label for="roleId">${$('#strRole').text()} ID</label>` +
        `<input id="roleId" class="form-control" type="number" name="id" value="${id || ''}" ${id ? 'readonly' : ''}>` +
        '</div>' +
        '<div class="form-group">' +
        `<label for="name">${$('#strName').text()}</label>` +
        `<input id="name" class="form-control" type="text" name="name" value="${name || ''}">` +
        '</div>' +
        '</form>';
    return html;
}

function saveFunctions() {
    const roleId = parseInt($('#roleSelect').val());
    let data = {id: roleId};

    data.functions = $('#functionData input:checkbox:checked').map(function () {
        return {
            roleId,
            functionId: parseInt($(this).val())
        }
    }).get();
    $.ajax({
        url: '/role/functions',
        method: 'PUT',
        dataType: 'json',
        data: JSON.stringify(data)
    }).done((r) => {
        BootstrapDialog.alert($('#strSaved').text());
    }).fail(alertErr);
}

function checkAll(dom) {
    const checked = $(dom).prop('checked');
    if (checked) {
        $('input:checkbox').prop('checked', true);
    } else {
        $('input:checkbox').prop('checked', false);
    }
}

function checkFunc(dom, id) {
    const checked = $(dom).prop('checked');
    let checkboxes;
    if (checked) {
        checkboxes = $(`.parent_${id}`).find('input:checkbox:not(:checked)')
    } else {
        checkboxes = $(`.parent_${id}`).find('input:checkbox:checked')
    }
    checkboxes.trigger('click');
    setCheckAllCB();
}

function setCheckAllCB() {
    const uncheckedLength = $('tbody input:checkbox:not(:checked)').length;
    $('#checkAll').prop('checked', uncheckedLength === 0);
}

function deleteRole() {
    const roleId = parseInt($('#roleSelect').val());
    BootstrapDialog.show({
        title: $('#strPrompt').text(),
        message: $('#strDeleteMsg').text(),
        buttons: [
            {
                label: $('#strBtnConfirm').text(),
                action: function (dialogItself) {
                    $.ajax({
                        url: `/role?id=${roleId}`,
                        method: 'delete',
                        dataType: 'json',
                    }).done(() => {
                        $('#roleSelect').find(':selected').remove();
                        $('#roleSelect').trigger('change');
                    }).fail(alertErr).always(() => dialogItself.close());
                }
            },
            {
                label: $('#strBtnCancel').text(),
                action: function (dialogItself) {
                    dialogItself.close();
                }
            }]
    });

}