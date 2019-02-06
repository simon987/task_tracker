import {MatPaginatorIntl} from "@angular/material";
import {TranslateService} from "@ngx-translate/core";

export function TranslatedPaginator(translate: TranslateService) {

    const paginatorIntl = new MatPaginatorIntl();

    getTranslations(translate, paginatorIntl);

    translate.onLangChange.subscribe(() => {
        getTranslations(translate, paginatorIntl)
    });

    return paginatorIntl;
}

function getTranslations(tr: TranslateService, p: MatPaginatorIntl) {

    tr.get("logs.first_page_label").subscribe((t) => p.firstPageLabel = t);
    tr.get("logs.last_page_label").subscribe((t) => p.lastPageLabel = t);
    tr.get("logs.items_per_page").subscribe((t) => p.itemsPerPageLabel = t);
    tr.get("logs.next_page").subscribe((t) => p.nextPageLabel = t);
    tr.get("logs.prev_page").subscribe((t) => p.previousPageLabel = t);
    tr.get("logs.of").subscribe((of) =>
        p.getRangeLabel = (page, pageSize, length) => `${page * pageSize + 1}-${Math.min(pageSize * (page + 1), length)} ${of} ${length}`
    );

}
