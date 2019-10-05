<template>
    <div class="card">
        <div class="card-body">
            <h3 class="card-title">Cluster status: {{ name }}</h3>

            <div class="card-text">
                <div class="alert alert-warning" role="alert" v-if="!isSameCluster">
                    Warning: This is not from the same cluster the browser got!
                </div>
                <div class="alert alert-danger" role="alert" v-if="!isCorrectCluster">
                    Error: This is the wrong cluster!
                </div>
                <div v-if="isOk">
    <div class="number" v-bind:class="{ numberNice: latency < 50 }">{{ latency }} ms</div>
    <table class="table">
        <thead>
        <tr>
            <th>Attribute</th>
            <th>Value</th>
        </tr>
        </thead>
        <tbody>
        <tr>
            <td>Cluster name</td>
            <td>{{ response.name }}</td>
        </tr>
        <tr>
            <td>Ingress to server HTTP version</td>
            <td>{{ response.proto }}</td>
        </tr>
        <tr>
            <td>HTTP endpoint to server</td>
            <td>{{ response.remote }}</td>
        </tr>
        <tr>
            <td>Response timestamp</td>
            <td>{{ response.time }}</td>
        </tr>
        <tr>
            <td>Client to Router ALPN</td>
            <td>{{ response.fwd_alpn }}</td>
        </tr>
        <tr>
            <td>Router ID</td>
            <td>{{ response.fwd_edge }}</td>
        </tr>
        <tr>
            <td>Client to Router port</td>
            <td>{{ response.fwd_port }}</td>
        </tr>
        <tr>
            <td>Client to Router HTTP version</td>
            <td>{{ response.fwd_proto }}</td>
        </tr>
        <tr>
            <td>Real client IP</td>
            <td>{{ response.fwd_remote }}</td>
        </tr>
        <tr>
            <td>Client to Router TLS 1.3 0RTT?</td>
            <td>{{ response.fwd_tls13early }}</td>
        </tr>
        <tr>
            <td>Client to Router TLS Cipher</td>
            <td>{{ response.fwd_tlscipher }}</td>
        </tr>
        <tr>
            <td>Client to Router TLS version</td>
            <td>{{ response.fwd_tlsver }}</td>
        </tr>
        <tr>
            <td>Client to Router rotocol</td>
            <td>{{ response.fwd_uri }}</td>
        </tr>
        <tr>
            <td>Kubernetes version</td>
            <td>{{ response.version }}</td>
        </tr>
        <tr>
            <td>Kubernetes master URL</td>
            <td>{{ response.master }}</td>
        </tr>
        </tbody>
    </table>
                </div>
            <div class="alert alert-danger" role="alert" v-else>
                Couldn't connect!
            </div>
    </div>
    </div>
    </div>
</template>

<script>
export default {
    name: "Ping",
    props: {
        remote: String,
        name: String,
    },
    data: function () {
        return {
            latency: 0,
            pre: undefined,
            timer: undefined,
            previous: undefined,
            response: {},
            isSameCluster: false,
            isCorrectCluster: false,
            isOk: true,
        }
    },
    methods: {
        handleLatencyResponse(response) {
            var post = new Date();
            this.latency = post.getUTCMilliseconds() - this.pre.getUTCMilliseconds();

            // Let's see whether we get data from _the same cluster_ as the browser originally got on the first page load
            if (typeof this.previous === "undefined") {
                // First invocation -> get values from window object
                this.previous = window.clusters[this.$props.name];
            }

            var expression = new RegExp("window.clusters[^{]+", "g");
            var strResp = response.data.replace(expression, "");
            this.response = JSON.parse(strResp);

            if (typeof this.previous === "undefined") {
                this.isSameCluster = false;
            } else {
                this.isSameCluster = this.response.name === this.previous.name;
            }
            this.isCorrectCluster = this.response.name === this.$props.name;
            this.isOk = true;
        },
        handleErrorResponse(error) {
            this.isOk = false;
            // eslint-disable-next-line
            console.log(error);
        },
        updateLatency() {
            this.pre = new Date();
            this.$http.get(this.$props.remote + "/ping").then(this.handleLatencyResponse).catch(this.handleErrorResponse)
        }
    },
    created() {
        this.timer = setInterval(this.updateLatency, 1000)
    },
    beforeDestroy() {
        clearInterval(this.timer)
    }
}
</script>

<style scoped>
.number {
    color: darkorange;
    font-size: 3rem;
    padding: 2rem;
    padding-top: 0.2rem;
}
    .numberNice {
        color: green;
    }
</style>