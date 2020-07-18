var vm = new Vue({
    el: "#App",

    data: function () {
        return {
            loading_flag: false,
            List: [],
            postData: {
                cmdJob: "",
                time: {
                    second: "*",
                    minute: "*",
                    hour: "*",
                    day: "*",
                    month: "*",
                    dow: "*"
                }
            }
        };
    },

    computed: {
        jobTime() {
            return `${this.postData.time.second} ${this.postData.time.minute} ${this.postData.time.hour} ${this.postData.time.day} ${this.postData.time.month} ${this.postData.time.dow}`
        }
    },

    mounted: function () {
        this.load();
    },

    methods: {
        load () {
            return fetch("/status")
                .then((response) => {
                    return response.json();
                })
                .then((resp) => {
                    this.List = resp;
                    return Promise.resolve();
                });
        },

        remove(id) {
            if (confirm("Are you sure ?")) {
                var data = new FormData();
                this.loading_flag = true;
                data.append("id", id);

                fetch("/remove", {
                        method: "POST",
                        body: data
                    })
                    .then((response) => {
                        return response.json();
                    })
                    .then((resp) => {
                        if (resp.status) {
                            return this.load();
                        }
                        return Promise.resolve();
                    })
                    .finally(() => {
                        this.loading_flag = false;
                    });
            }
        },

        add () {
            var data = new FormData();
            data.append("time", this.jobTime);
            data.append("cmd", this.postData.cmdJob);

            this.loading_flag = true;

            fetch("/add", {
                    method: "POST",
                    body: data
                })
                .then((response) => {
                    return response.json();
                })
                .then((resp) => {
                    if (resp.status) {
                        this.resetForm();
                        return this.load();
                    }
                    return Promise.resolve();
                })
                .finally(() => {
                    this.loading_flag = false;
                });
        },

        resetForm() {
            this.postData = {
                cmdJob: "",
                time: {
                    second: "*",
                    minute: "*",
                    hour: "*",
                    day: "*",
                    month: "*",
                    dow: "*"
                }
            };
        }
    }
});