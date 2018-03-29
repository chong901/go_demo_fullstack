let parameters;
let units;
$(() => {
    $.ajax({
        url: '/api/v1/configuration/prodParameters',
    }).done((r) => {
        parameters = r.data || [];
    }).fail(alertErr);

    $('#form').on('submit', function(ev){
        ev.preventDefault();

        let data = getObjFromSerialize($(this).serializeArray());
        data.parameters = $('.paramsTr').map(function () {
            const $this = $(this);
            return {
                id: parseInt($this.find('.id').val()),
                key: $this.find('.pKey').val(),
                value: $this.find('.pValue').val(),
                skuId: $('#skuId').val()
            }
        }).get();

        $.ajax({
            url: '/skuM/save',
            method: 'post',
            data: JSON.stringify(data)
        }).done((r) => {
            let formData = new FormData($('#form')[0]);
            const id = r.id;
            formData.append("id", id);

            $.ajax({
                url: '/skuM/uploadFile',
                method: 'post',
                contentType: false,
                processData: false,
                data: formData
            }).done(() => {
                window.location.assign(`/skuM/save?id=${id}`);
            }).fail(alertErr);
        }).fail(alertErr);
    })
});

function addCond() {
    let html = '<div class="tr dataRow paramsTr">';
    html += '<span class="td">' + createPamamsSelect(parameters) + '</span>';
    html += '<span class="td"><input class="form-control tdInput pValue"></span>';
    html += `<span class="td"><input type="button" class="btn btn-danger condBtnRm" value="${$('#strDelete').text()}" onclick="removeParam(this)"></span>`;
    html += '</div>';
    $('#condTable').append(html);
}

function createPamamsSelect(parameters) {
    return parameters.reduce((pre, cur) => {
        if (cur.Name && cur.Name !== 0) {
            return pre + `<option value="${cur.Name}">${cur.Name}</option>`;
        } else {
            return pre;
        }
    }, `<select class="form-control tdInput pKey"><option></option>`) + '</select>';
}
