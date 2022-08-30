/*global browser, chrome*/
const browser = browser || chrome;

const host = "http://127.0.0.1:8080";

function sendToTagsrus(q) {
  const url = new URL(host);
  url.pathname = "/api/import";
  url.search = new URLSearchParams({
    link: encodeURIComponent(q),
  }).toString();

  fetch(url);
}

browser.contextMenus.create({
  title: "Send to Tagsrus",
  contexts: ["link"],
  onclick: (item) => {
    sendToTagsrus(item.linkUrl);
  },
});

browser.browserAction.onClicked.addListener((tab) => {
  sendToTagsrus(tab.url);
});