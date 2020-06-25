const fs = require('fs');
const xml2js = require('xml2js');
const parser = new xml2js.Parser();
const database = require('./database');

async function loadDocData(dataId) {
    const data = await database.getContent(dataId);

    parser.parseString(data.toString(), function (err, result) {
        doc = result.doc
        title = doc.title[0]
        url = doc.url[0]
        abstract = doc.abstract[0];
        sections = doc.links[0].sublink.map((link) => link.anchor[0]);
    });

    return {
        title: title,
        url: url,
        abstract: abstract,
        sections: sections
    }
}

async function getQueryPrices(dataId, queries) {
    var doc = await loadDocData(dataId)

    var prices = {
        titlePrice: doc.title.length,
        urlPrice: doc.url.length,
        abstractPrice: doc.abstract.length,
        sectionsPrice: doc.sections.length * 10
    }

    var priceArray = '['
    for (i in queries) {
        if (queries[i] === 'Title')
            priceArray += prices.titlePrice
        else if (queries[i] === 'Url')
            priceArray += prices.urlPrice
        else if (queries[i] === 'Abstract')
            priceArray += prices.abstractPrice
        else if (queries[i] === 'Sections')
            priceArray += prices.sectionsPrice
        else
            priceArray += 99999.0 // query invalida, preço adicionado para que no chaincode seja lançado erro de query invalida 
        if (i != queries.length - 1) priceArray += ','
    }
    priceArray += ']'
    return priceArray
}

async function getResponseContent(dataId, query, queryPrice, pricePaid) {
    var doc = await loadDocData(dataId);
    var field;

    const percentage = pricePaid / queryPrice;

    if (query === 'Title')
        field = doc.title;
    else if (query === 'Url')
        field = doc.url;
    else if (query === 'Abstract')
        field = doc.abstract;
    else if (query === 'Sections')
        field = doc.sections;

    return field.slice(0, field.length * percentage).toString();
}

module.exports = {
    getQueryPrices: getQueryPrices,
    getResponseContent: getResponseContent
}