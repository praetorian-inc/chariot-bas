# Chariot BAS

This is the home to all Tactics, Techniques and Procedures (TTP) for Chariot's internal assessments.

What makes a good TTP? Code that executes a known adversarial behavior expected to fail in a secure environment.

## Usage

Using a TTP in Chariot involves:

1. Build the binary
2. Provision an agent
3. Upload the binary
4. Schedule the test
5. Execute the test

### Build

Build any TTP using this command:

```bash
# enter TTP uuid
uuid=<uuid>
# supported: linux, darwin and windows
platform="linux"
# ensure unique hash
echo "// $(date)" >> tests/<uuid>.go
# build the binary
GOOS=$platform GOARCH=amd64 go build -o "${uuid}-${platform}" tests/<uuid>-<platform>.go
```

### Provision

Provision a new agent using the following command:

```bash
praetorian chariot add asset --name <agentid>
```

### Upload

Use the Praetorian CLI to upload the binary to your account:

```bash
praetorian chariot add file malware/${uuid}-${platform}
```

### Schedule

Use the Praetorian CLI to schedule the task:

```bash
praetorian chariot add job <uuid> -asset <agentid>
```

### Execute

Deploy the agent to any endpoint. It will check for scheduled tests every 60 seconds. The following curl commands can be used to download a copy of the agent for your account.

Agent.sh (Darwin and Linux):
```bash
curl -o agent.sh https://d0qcl2e18h.execute-api.us-east-2.amazonaws.com/chariot/<account>/<agentid>/linux
```

Agent.ps1 (Windows):
```bash
curl -o agent.ps1 https://d0qcl2e18h.execute-api.us-east-2.amazonaws.com/chariot/<account>/<agentid>/windows
```

## Contributing

To write a TTP, create a ``.go`` file in the ``internal/tests`` directory, using the template below.

```go
package tests

import "github.com/praetorian-inc/chariot-bas/internal/endpoint"

func test() {
    // STOP with a predefined condition
    // review codes.go for all options
    endpoint.Stop(endpoint.Risk.Allowed)
}

func cleanup() {
    // optional logic to reverse the effects of this test
}

func main(){
    endpoint.Start(test, cleanup)
}
```

## Endpoint SDK

Each Chariot TTP uses the endpoint module, an SDK for common operations. Review the options in ``internal/endpoint``.
