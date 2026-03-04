async function nextEpisode(id) {
    const url = '/update?id=' + id;
    const response = await fetch(url, { method: "POST" })
    location.reload()
}

async function prevEpisode(id) {
    const url = '/downdate?id=' + id;
    await fetch(url, { method: "POST" })
    location.reload()
}