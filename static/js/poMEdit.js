let skus;
let units;
$(function () {
    $.when(
        $.ajax('/api/v1/sku'),
        $.ajax('/api/v1/configuration/units')
    ).done((r1, r2) => {
        skus = r1[0].data || [];
        units = r2[0].data || [];
    }).fail(alertErr);

    $('#form').on('submit', function (ev) {
        ev.preventDefault();

        let data = getObjFromSerialize($(this).serializeArray(), null, ['due']);
        data.skuList = getSkuData() || [];

        $.ajax({
            url: '/poM/save',
            method: 'post',
            dataType: 'json',
            data: JSON.stringify(data)
        }).done((r) => {
            window.location.assign(`/poM/save?id=${r.id}`);
        }).fail(alertErr);

    })
});

function getSkuData() {
    return $('.dataRow').map(function () {
        const $tr = $(this);
        return {
            id: parseInt($tr.find('.id').val()),
            skuId: $tr.find('.skuId').val(),
            quantity: parseInt($tr.find('.quantity').val()),
            unit: $tr.find('.unit').val(),
            poId: $('#poId').val()
        }
    }).get();
}

function createSkuSelect(skuArr) {
    let html = '<select class="form-control trSelect tdInput skuId">';
    html += skuArr.reduce(function (pre, cur) {
        let temp = '<option value="' + cur.id + '"';
        temp += '>';
        temp += cur.id;
        temp += '</option>';
        return pre + temp;
    }, '<option></option>');

    html += '</select>';
    return html;
}

function createUnitSelect(units) {
    return units.reduce((pre, cur) => {
        if (!cur) return pre;
        return pre + `<option value="${cur.Name}">${cur.Name}</option>`;
    }, '<select class="form-control trSelect tdInput unit"><option></option>') + '</select>';
}

function createEmptySku(skuArr) {
    let html = '<div class="tr dataRow">';
    html += '<span class="td">' + createSkuSelect(skuArr) + '</span>';
    html += '<span class="td"><input type="number" class="form-control tdInput quantity"></span>';
    html += `<span class="td">${createUnitSelect(units)}</span>`;
    html += '<span class="td">';
    html += `<input type="button" class="btn btn-danger" value="${$('#strDelete').text()}" onclick="removePoSku(this)">`;
    html += '</span>';
    html += '</div>';
    return html;
}

function addSku(dom) {
    let table = $(dom).siblings('.table');
    const length = table.find('.dataRow').length;
    table.append(createEmptySku(skus, length));
}
