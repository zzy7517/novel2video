package llm

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

const (
	projectID = "p-sgx9cobj"
	location  = "asia-east1"
)

var geminiClient *genai.Client

func InitClient() {
	var err error
	options := getGrpcOptions()
	geminiClient, err = genai.NewClient(context.Background(), projectID, location, options...)
	if err != nil {
		log.Fatalf("error creating gemini client: %+v,location:%s", err, location)
	}
}

func GetClient() *genai.Client {
	if geminiClient == nil {
		InitClient()
	}
	return geminiClient
}

func GetSafetySetting() []*genai.SafetySetting {
	return []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryUnspecified,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
	}
}

func QueryGemini(ctx context.Context, input, systemContent string, modelName string, temperature float32, maxOutputTokens int32) (res string, err error) {
	client := GetClient()
	if client == nil {
		return res, fmt.Errorf("error getting gemini client")
	}
	gemini := client.GenerativeModel(modelName)
	gemini.SetTemperature(temperature)
	if maxOutputTokens != 0 {
		gemini.SetMaxOutputTokens(maxOutputTokens)
	}
	gemini.SafetySettings = GetSafetySetting()
	if len(systemContent) > 0 {
		var sysPart []genai.Part
		sysPart = append(sysPart, genai.Text(systemContent))
		gemini.SystemInstruction = &genai.Content{Parts: sysPart}
	}
	var prompt []genai.Part
	prompt = append(prompt, genai.Text(input))
	resp, err := gemini.GenerateContent(ctx, prompt...)
	if err != nil {
		return res, err
	}
	if len(resp.Candidates) == 0 {
		return res, fmt.Errorf("get result from gemini failed, nil response")
	}
	if resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return res, fmt.Errorf("get result from gemini failed, nil response")
	}
	t, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if ok {
		return string(t), nil
	}
	return res, fmt.Errorf("get result from gemini failed")
}

func getGrpcOptions() []option.ClientOption {
	credentialsFile := "/Users/zhongyuanzhang/Desktop/gemini-credential.json"
	return []option.ClientOption{
		option.WithCredentialsFile(credentialsFile),
		option.WithGRPCDialOption(grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)),
		// getGrpcProxy(),
		// option.WithGRPCConnectionPool(500),
	}
}

func getGrpcProxy() option.ClientOption {
	f := WithProxyDialer("http://proxy.p1.cn:1338")
	return option.WithGRPCDialOption(
		grpc.WithContextDialer(
			func(ctx context.Context, addr string) (net.Conn, error) {
				if deadline, ok := ctx.Deadline(); ok {
					return f(addr, time.Until(deadline))
				}
				return f(addr, 0)
			}),
	)
}

func WithProxyDialer(proxyHost string) func(addr string, duration time.Duration) (net.Conn, error) {
	return func(addr string, duration time.Duration) (net.Conn, error) {
		network, proxyAddr := parseDialTarget(proxyHost)
		conn, err := dialContext(network, proxyAddr)
		if err != nil {
			return nil, err
		}
		return doHTTPConnectHandshake(conn, addr)
	}
}

func parseDialTarget(target string) (net string, addr string) {
	net = "tcp"

	m1 := strings.Index(target, ":")
	m2 := strings.Index(target, ":/")

	// handle unix:addr which will fail with url.Parse
	if m1 >= 0 && m2 < 0 {
		if n := target[0:m1]; n == "unix" {
			net = n
			addr = target[m1+1:]
			return net, addr
		}
	}
	if m2 >= 0 {
		t, err := url.Parse(target)
		if err != nil {
			return net, target
		}
		scheme := t.Scheme
		addr = t.Path
		if scheme == "unix" {
			net = scheme
			if addr == "" {
				addr = t.Host
			}
			return net, addr
		}
	}

	return net, target
}

func dialContext(network, address string) (net.Conn, error) {
	return (&net.Dialer{}).Dial(network, address)
}

type bufConn struct {
	net.Conn
	r io.Reader
}

func (c *bufConn) Read(b []byte) (int, error) {
	return c.r.Read(b)
}

func doHTTPConnectHandshake(conn net.Conn, addr string) (_ net.Conn, err error) {
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	req := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Host: addr},
		Header: map[string][]string{"User-Agent": {"grpc-go/1.11.3"}},
	}

	if err := sendHTTPRequest(req, conn); err != nil {
		return nil, fmt.Errorf("failed to write the HTTP request: %v", err)
	}

	r := bufio.NewReader(conn)
	resp, err := http.ReadResponse(r, req)
	if err != nil {
		return nil, fmt.Errorf("reading server HTTP response: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, fmt.Errorf("failed to do connect handshake, status code: %s", resp.Status)
		}
		return nil, fmt.Errorf("failed to do connect handshake, response: %q", dump)
	}

	return &bufConn{Conn: conn, r: r}, nil
}

func sendHTTPRequest(req *http.Request, conn net.Conn) error {
	if err := req.Write(conn); err != nil {
		return fmt.Errorf("failed to write the HTTP request: %v", err)
	}
	return nil
}
