function addConfig(t) {
    BootstrapDialog.show({
        size: BootstrapDialog.SIZE_SMALL,
        title: $('#strAdd').text(),
        message: createAddMessage(),
        buttons: [
            {
                label: $('#strBtnConfirm').text(),
                action: function (bd) {
                    $.ajax({
                        url: '/configuration',
                        method: 'post',
                        data: {
                            type: t,
                            name: $('#name').val()
                        }
                    }).done((r) => {
                        $('#emptyTr').remove();
                        $('#configTable').append(createTr(r.id, r.name));
                        bd.close();
                    }).fail(alertErr);
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

function createAddMessage() {
    return '<div class="form-group">' +
        `<label for="name">${$('.strName').text()}</label>` +
        '<input class="form-control" id="name" >' +
        '</div>';
}

function createTr(id, name) {
    return '<tr>' +
        `<td>${name}</td>` +
        `<td><input class="btn btn-danger" type="button" value="${$('#strDelete').text()}" onclick="remove(this, '/configuration', 'id', '${id}')"></td>` +
        '</tr>';
}