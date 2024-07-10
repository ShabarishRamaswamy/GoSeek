window.addEventListener("load", async () => {
    console.log("Loaded !\n", "Doing speed test now !\n");
    speedTestTime = await doSpeedTest();
    speedTestPong(speedTestTime);
});

async function doSpeedTest() {
    var collector = [];

    for (var i = 0; i < 5; i++) {
        var value = await speedTestPing();
        collector.push(value);
    }

    // Deal with the outliers.
    // Usually the first request gives unexpectedly large time readings.
    var finalSpeed = Infinity;
    for (var i = 0; i < 5; i++) {
        var currentAvg = 0;
        for (var j = 0; j < 5; j++) {
            if (j == i) {
                continue;
            }
            currentAvg += collector[j];
        }
        currentAvg = currentAvg / 5;
        // console.log("Average Number: ", i, " Which is: ", currentAvg);
        if (currentAvg < finalSpeed) {
            finalSpeed = currentAvg;
        }
    }

    // console.log("Final Speed: ", finalSpeed);
    return finalSpeed;
}

async function speedTestPing() {
    try {
        // TODO: Optimize this, currently I spend 10ms in the time calculation.
        // For the sake of this application, we will assume 90% of the time calculation is accurate.
        // This is a Heuristic only.
        var startTime = new Date().getTime();
        const resp = await fetch("/speedTest/request");
        var endTime = new Date().getTime();

        // const text = await resp.text();
        console.log("Time taken: ", endTime - startTime);
        return endTime - startTime;
    } catch (err) {
        console.error("Error: ", err);
        return 0;
    }
}

async function speedTestPong(speedTestTime) {
    try {
        await fetch("/speedTest/response", {
            body: JSON.stringify({ time: speedTestTime }),
            method: "POST",
        });
    } catch (e) {
        console.log(e);
    }
}
