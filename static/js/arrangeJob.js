let noDataTr = '<tr class="noData text-center"><td>No Data</td></tr>';
const sortableJQSelect = "tr:not(.noData):not(.processJob)";
const PLANNED_TYPE = 1;
const UNDONE_TYPE = 2;
let jobTable;

const iconMinus = $('<a class="icon-minus" href="#"><span class="glyphicon glyphicon-minus"></span></a>').on('click', clickIconMinus);
const iconPlus = $('<a class="icon-plus" href="#"><span class="glyphicon glyphicon-plus"></span></a>').on('click', clickIconPlus);

$(() => {
    jobTable = $('.hiddenTable');
    setJobTables();

    const undoneTable = $('#undone .sortable');

    $('#planned .sortable').sortable({
        items: sortableJQSelect,
        receive: function (e, ui) {
            checkSenderTrEmpty(ui);

            ui.item.siblings('.noData').remove();
            ui.item.before(ui.item.data('items'));

            $(this).find('.icon-plus').replaceWith(iconMinus.clone(true));
            $(this).find('tr').unbind();
        }
    });

    $('#undone .sortable').sortable({
        connectWith: '#planned .sortable',
        items: sortableJQSelect,
        helper: multipleSortableHelper,
        start: function (e, ui) {
            let elements = ui.item.siblings('.active').not('.ui-sortable-placeholder');
            ui.item.data('items', elements);
        },
        receive: function (e, ui) {
            checkSenderTrEmpty(ui);

            ui.item.siblings('.noData').remove();

            ui.item.on('click', function(e) {
                $(this).toggleClass('active');
            });
        },
        stop: function (e, ui) {
            ui.item.siblings('.active').removeClass('hidden');
            $('table .active').removeClass('active');
        }
    });

    $('.btnMachine').on('click', function(){
        $(this).addClass('active');
        $(this).siblings('.active').removeClass('active');
        const machineId = $(this).data('machine-id');
        let ajaxJobPlanned = $.ajax({
            url: `/api/v1/jobPlanned?machineId=${machineId}`
        });
        let ajaxJobUndone = $.ajax({
           url: '/api/v1/jobUndone'
        });

        $.when(ajaxJobPlanned, ajaxJobUndone).done((jobPlanned, jobUndone) => {
            setJobPlannedTbody(jobPlanned[0]);
            setJobUndoneTbody(jobUndone[0]);
        }).fail(alertErr);
    });

    $('#planned .btnSave').on('click', function(){
        const currentMachine = parseInt($('.btnMachine.active').data('machine-id'));
        let plannedJob = $('#planned .sortable>tr');
        let unplannedJob = $('#undone .sortable>tr');

        let plannedJobData = [];
        let unplannedJobData;
        let nextJobId = '';

        plannedJob.get().reverse().forEach((k, i) => {
            const currentId = $(k).data('id');
            const tempData = {
                id: currentId,
                nextJobId: nextJobId,
                machineId: currentMachine,
                status: $(k).data('status')
            };
            plannedJobData.push(tempData);
            nextJobId = currentId;
        });

        plannedJobData = plannedJobData.reverse().map((k, i) => {
            k.jobOrder = i + 1;
            return k;
        });

        unplannedJobData = unplannedJob.map((i, k) => {
            return $(k).find('.id').html();
        }).get();

        const ajaxData = {
            plannedJobs: plannedJobData,
            unplannedJobIds: unplannedJobData
        };

        $.ajax({
            dataType: 'json',
            url: '/jobSchedule',
            method: 'put',
            data: JSON.stringify(ajaxData)
        }).done(() => {
            BootstrapDialog.alert($('#strSaved').text());
        }).fail(alertErr);
    });

    $('.btnMachine').first().click();

    if(undoneTable.find('tr').length === 0){
        addNoDataTr(undoneTable);
    }
});

function multipleSortableHelper(e, item) {
    if(!item.hasClass('active'))
        item.addClass('active');
    let elements = $(this).find('.active').not('.ui-sortable-placeholder').clone();
    let helper = $('<tbody/>');
    item.siblings('.active').addClass('hidden');
    return helper.append(elements);
}

function checkSenderTrEmpty(ui){
    let sender = ui.sender;
    if(sender.find('tr').not('.active').length == 0){
        addNoDataTr(sender);
    }
}

function addNoDataTr(jqObj){
    const colspanNum = jqObj.closest('table').find('th').length;
    let noDataTrNew = $(noDataTr).clone();
    noDataTrNew.find('td').prop('colspan', colspanNum);
    jqObj.html(noDataTrNew);
}

function createJobTr(data, type){
    if(!data){
        return '';
    }

    let idAppend = '';
    let trCls = '';
    let output = getOutputByJob(data.id);
    switch (data.status){
        case 4:
            idAppend = ' (Processing)';
            trCls = 'processJob';
    }
    let tr = $('<tr></tr>');
    if(trCls){
        tr.addClass(trCls);
    }
    tr.data('status', data.status);
    tr.data('id', data.id);
    tr.append(`<td class="id">${data.id + idAppend}</td>`)
        .append(`<td>${data.skuId}</td>`)
        .append(`<td>${data.quantity}</td>`)
        .append(`<td>${output}</td>`)
        .append($('<td></td>').append(getIcon(type)));
    return tr;
}

function getIcon(type){
    if(!type){
        return '';
    }
    switch(type){
        case PLANNED_TYPE:
            return iconMinus.clone(true);
        case UNDONE_TYPE:
            return iconPlus.clone(true);
    }
}

function clickIconPlus(){
    const tr = $(this).closest('tr');
    tr.find('a').replaceWith(iconMinus.clone(true));
    tr.removeClass('active');
    tr.unbind();
    $('#planned .sortable').append(tr);
    checkBothTable();
}

function clickIconMinus(e){
    e.stopPropagation();

    const tr = $(this).closest('tr');
    tr.find('a').replaceWith(iconPlus.clone(true));
    tr.bind('click', function(e){
        $(this).toggleClass('active');
    });
    $('#undone .sortable').append(tr);
    checkBothTable();
}

function setJobPlannedTbody(data){
    const plannedTBody = $('#planned .sortable');
    if(!data || data.totalLength === 0){
        addNoDataTr(plannedTBody);
        return
    }
    let processJobs = data.processJobs || [];
    let plannedJobs = data.plannedJobs || [];

    processJobs = processJobs.map((d) => createJobTr(d));
    plannedJobs = plannedJobs.map((d) => createJobTr(d, PLANNED_TYPE));

    plannedTBody.html(processJobs.concat(plannedJobs));
}

function setJobUndoneTbody(data){
    const undoneTbody = $('#undone .sortable');
    if(!data || !data.undoneJobs || data.undoneJobs.length === 0){
        addNoDataTr(undoneTbody);
        return
    }

    let undoneJobs = data.undoneJobs.map((d) => createJobTr(d, UNDONE_TYPE).on('click', function(e){
        $(this).toggleClass('active');
    }));

    undoneTbody.html(undoneJobs);
}

function setJobTables(){
    $('#planned>.top').html(jobTable.clone());
    $('#undone>.content').html(jobTable.clone());
}

function checkBothTable(){
    const plannedTable = $('#planned .sortable');
    const undoneTable = $('#undone .sortable');

    checkTableEmpty(plannedTable);
    checkTableEmpty(undoneTable);
}

function checkTableEmpty(jqObj){
    if(jqObj.find('tr').length === 0){
        addNoDataTr(jqObj)
    }else{
        jqObj.find('.noData').remove();
    }
}

function getOutputByJob(jobId){
    let pcs = 0

    return pcs
}
