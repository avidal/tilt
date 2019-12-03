// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cli

import (
	"context"
	"time"

	"github.com/google/wire"
	"github.com/jonboulle/clockwork"
	"github.com/windmilleng/wmclient/pkg/dirs"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/windmilleng/tilt/internal/analytics"
	"github.com/windmilleng/tilt/internal/build"
	"github.com/windmilleng/tilt/internal/cloud"
	"github.com/windmilleng/tilt/internal/cloud/cloudurl"
	"github.com/windmilleng/tilt/internal/container"
	"github.com/windmilleng/tilt/internal/containerupdate"
	"github.com/windmilleng/tilt/internal/docker"
	"github.com/windmilleng/tilt/internal/dockercompose"
	"github.com/windmilleng/tilt/internal/dockerfile"
	"github.com/windmilleng/tilt/internal/engine"
	analytics2 "github.com/windmilleng/tilt/internal/engine/analytics"
	"github.com/windmilleng/tilt/internal/engine/configs"
	"github.com/windmilleng/tilt/internal/engine/dockerprune"
	"github.com/windmilleng/tilt/internal/engine/k8swatch"
	"github.com/windmilleng/tilt/internal/engine/runtimelog"
	"github.com/windmilleng/tilt/internal/feature"
	"github.com/windmilleng/tilt/internal/hud"
	"github.com/windmilleng/tilt/internal/hud/server"
	"github.com/windmilleng/tilt/internal/k8s"
	"github.com/windmilleng/tilt/internal/minikube"
	"github.com/windmilleng/tilt/internal/store"
	"github.com/windmilleng/tilt/internal/synclet/sidecar"
	"github.com/windmilleng/tilt/internal/tiltfile"
	"github.com/windmilleng/tilt/internal/tiltfile/k8scontext"
	"github.com/windmilleng/tilt/internal/token"
	"github.com/windmilleng/tilt/pkg/model"
)

// Injectors from wire.go:

func wireDockerPrune(ctx context.Context, analytics2 *analytics.TiltAnalytics) (dpDeps, error) {
	clientConfig := k8s.ProvideClientConfig()
	config, err := k8s.ProvideKubeConfig(clientConfig)
	if err != nil {
		return dpDeps{}, err
	}
	env := k8s.ProvideEnv(ctx, config)
	restConfigOrError := k8s.ProvideRESTConfig(clientConfig)
	clientsetOrError := k8s.ProvideClientset(restConfigOrError)
	portForwardClient := k8s.ProvidePortForwardClient(restConfigOrError, clientsetOrError)
	namespace := k8s.ProvideConfigNamespace(clientConfig)
	kubeContext, err := k8s.ProvideKubeContext(config)
	if err != nil {
		return dpDeps{}, err
	}
	int2 := provideKubectlLogLevel()
	kubectlRunner := k8s.ProvideKubectlRunner(kubeContext, int2)
	client := k8s.ProvideK8sClient(ctx, env, restConfigOrError, clientsetOrError, portForwardClient, namespace, kubectlRunner, clientConfig)
	runtime := k8s.ProvideContainerRuntime(ctx, client)
	minikubeClient := minikube.ProvideMinikubeClient()
	clusterEnv, err := docker.ProvideClusterEnv(ctx, env, runtime, minikubeClient)
	if err != nil {
		return dpDeps{}, err
	}
	localEnv, err := docker.ProvideLocalEnv(ctx, clusterEnv)
	if err != nil {
		return dpDeps{}, err
	}
	localClient := docker.ProvideLocalCli(ctx, localEnv)
	clusterClient, err := docker.ProvideClusterCli(ctx, localEnv, clusterEnv, localClient)
	if err != nil {
		return dpDeps{}, err
	}
	switchCli := docker.ProvideSwitchCli(clusterClient, localClient)
	extension := k8scontext.NewExtension(kubeContext, env)
	dockerComposeClient := dockercompose.NewDockerComposeClient(localEnv)
	defaults := _wireDefaultsValue
	tiltfileLoader := tiltfile.ProvideTiltfileLoader(analytics2, client, extension, dockerComposeClient, defaults)
	cliDpDeps := newDPDeps(switchCli, tiltfileLoader)
	return cliDpDeps, nil
}

var (
	_wireDefaultsValue = feature.MainDefaults
)

func wireCmdUp(ctx context.Context, hudEnabled hud.HudEnabled, analytics3 *analytics.TiltAnalytics, cmdUpTags analytics2.CmdUpTags) (CmdUpDeps, error) {
	v := provideClock()
	renderer := hud.NewRenderer(v)
	modelWebHost := provideWebHost()
	modelWebPort := provideWebPort()
	webURL, err := provideWebURL(modelWebHost, modelWebPort)
	if err != nil {
		return CmdUpDeps{}, err
	}
	headsUpDisplay, err := hud.ProvideHud(hudEnabled, renderer, webURL, analytics3)
	if err != nil {
		return CmdUpDeps{}, err
	}
	reducer := _wireReducerValue
	storeLogActionsFlag := provideLogActions()
	storeStore := store.NewStore(reducer, storeLogActionsFlag)
	clientConfig := k8s.ProvideClientConfig()
	config, err := k8s.ProvideKubeConfig(clientConfig)
	if err != nil {
		return CmdUpDeps{}, err
	}
	env := k8s.ProvideEnv(ctx, config)
	restConfigOrError := k8s.ProvideRESTConfig(clientConfig)
	clientsetOrError := k8s.ProvideClientset(restConfigOrError)
	portForwardClient := k8s.ProvidePortForwardClient(restConfigOrError, clientsetOrError)
	namespace := k8s.ProvideConfigNamespace(clientConfig)
	kubeContext, err := k8s.ProvideKubeContext(config)
	if err != nil {
		return CmdUpDeps{}, err
	}
	int2 := provideKubectlLogLevel()
	kubectlRunner := k8s.ProvideKubectlRunner(kubeContext, int2)
	client := k8s.ProvideK8sClient(ctx, env, restConfigOrError, clientsetOrError, portForwardClient, namespace, kubectlRunner, clientConfig)
	ownerFetcher := k8s.ProvideOwnerFetcher(client)
	podWatcher := k8swatch.NewPodWatcher(client, ownerFetcher)
	nodeIP, err := k8s.DetectNodeIP(ctx, env)
	if err != nil {
		return CmdUpDeps{}, err
	}
	serviceWatcher := k8swatch.NewServiceWatcher(client, ownerFetcher, nodeIP)
	podLogManager := runtimelog.NewPodLogManager(client)
	portForwardController := engine.NewPortForwardController(client)
	fsWatcherMaker := engine.ProvideFsWatcherMaker()
	timerMaker := engine.ProvideTimerMaker()
	watchManager := engine.NewWatchManager(fsWatcherMaker, timerMaker)
	runtime := k8s.ProvideContainerRuntime(ctx, client)
	minikubeClient := minikube.ProvideMinikubeClient()
	clusterEnv, err := docker.ProvideClusterEnv(ctx, env, runtime, minikubeClient)
	if err != nil {
		return CmdUpDeps{}, err
	}
	localEnv, err := docker.ProvideLocalEnv(ctx, clusterEnv)
	if err != nil {
		return CmdUpDeps{}, err
	}
	localClient := docker.ProvideLocalCli(ctx, localEnv)
	clusterClient, err := docker.ProvideClusterCli(ctx, localEnv, clusterEnv, localClient)
	if err != nil {
		return CmdUpDeps{}, err
	}
	switchCli := docker.ProvideSwitchCli(clusterClient, localClient)
	dockerContainerUpdater := containerupdate.NewDockerContainerUpdater(switchCli)
	syncletImageRef, err := sidecar.ProvideSyncletImageRef(ctx)
	if err != nil {
		return CmdUpDeps{}, err
	}
	syncletManager := containerupdate.NewSyncletManager(client, syncletImageRef)
	syncletUpdater := containerupdate.NewSyncletUpdater(syncletManager)
	execUpdater := containerupdate.NewExecUpdater(client)
	engineUpdateModeFlag := provideUpdateModeFlag()
	updateMode, err := engine.ProvideUpdateMode(engineUpdateModeFlag, env, runtime)
	if err != nil {
		return CmdUpDeps{}, err
	}
	clock := build.ProvideClock()
	liveUpdateBuildAndDeployer := engine.NewLiveUpdateBuildAndDeployer(dockerContainerUpdater, syncletUpdater, execUpdater, updateMode, env, runtime, clock)
	labels := _wireLabelsValue
	dockerImageBuilder := build.NewDockerImageBuilder(switchCli, labels)
	imageBuilder := build.DefaultImageBuilder(dockerImageBuilder)
	cacheBuilder := build.NewCacheBuilder(switchCli)
	execCustomBuilder := build.NewExecCustomBuilder(switchCli, clock)
	clusterName := k8s.ProvideClusterName(ctx, config)
	kindPusher := engine.NewKINDPusher(env, clusterName)
	syncletContainer := sidecar.ProvideSyncletContainer(syncletImageRef)
	imageBuildAndDeployer := engine.NewImageBuildAndDeployer(imageBuilder, cacheBuilder, execCustomBuilder, client, env, analytics3, updateMode, clock, runtime, kindPusher, syncletContainer)
	dockerComposeClient := dockercompose.NewDockerComposeClient(localEnv)
	imageAndCacheBuilder := engine.NewImageAndCacheBuilder(imageBuilder, cacheBuilder, execCustomBuilder, updateMode)
	dockerComposeBuildAndDeployer := engine.NewDockerComposeBuildAndDeployer(dockerComposeClient, switchCli, imageAndCacheBuilder, clock)
	localTargetBuildAndDeployer := engine.NewLocalTargetBuildAndDeployer(clock)
	buildOrder := engine.DefaultBuildOrder(liveUpdateBuildAndDeployer, imageBuildAndDeployer, dockerComposeBuildAndDeployer, localTargetBuildAndDeployer, updateMode, env, runtime)
	compositeBuildAndDeployer := engine.NewCompositeBuildAndDeployer(buildOrder)
	buildController := engine.NewBuildController(compositeBuildAndDeployer)
	extension := k8scontext.NewExtension(kubeContext, env)
	defaults := _wireDefaultsValue
	tiltfileLoader := tiltfile.ProvideTiltfileLoader(analytics3, client, extension, dockerComposeClient, defaults)
	configsController := configs.NewConfigsController(tiltfileLoader, switchCli)
	dockerComposeEventWatcher := engine.NewDockerComposeEventWatcher(dockerComposeClient)
	dockerComposeLogManager := runtimelog.NewDockerComposeLogManager(dockerComposeClient)
	profilerManager := engine.NewProfilerManager()
	analyticsReporter := analytics2.ProvideAnalyticsReporter(analytics3, storeStore)
	tiltBuild := provideTiltInfo()
	webMode, err := provideWebMode(tiltBuild)
	if err != nil {
		return CmdUpDeps{}, err
	}
	webVersion := provideWebVersion(tiltBuild)
	assetsServer, err := provideAssetServer(webMode, webVersion)
	if err != nil {
		return CmdUpDeps{}, err
	}
	httpClient := cloud.ProvideHttpClient()
	address := cloudurl.ProvideAddress()
	snapshotUploader := cloud.NewSnapshotUploader(httpClient, address)
	headsUpServer, err := server.ProvideHeadsUpServer(ctx, storeStore, assetsServer, analytics3, snapshotUploader)
	if err != nil {
		return CmdUpDeps{}, err
	}
	modelNoBrowser := provideNoBrowserFlag()
	headsUpServerController := server.ProvideHeadsUpServerController(modelWebHost, modelWebPort, headsUpServer, assetsServer, webURL, modelNoBrowser)
	githubClientFactory := engine.NewGithubClientFactory()
	tiltVersionChecker := engine.NewTiltVersionChecker(githubClientFactory, timerMaker)
	analyticsUpdater := analytics2.NewAnalyticsUpdater(analytics3, cmdUpTags)
	eventWatchManager := k8swatch.NewEventWatchManager(client, ownerFetcher)
	cloudUsernameManager := cloud.NewUsernameManager(httpClient)
	updateUploader := cloud.NewUpdateUploader(httpClient, address)
	dockerPruner := dockerprune.NewDockerPruner(switchCli)
	v2 := engine.ProvideSubscribers(headsUpDisplay, podWatcher, serviceWatcher, podLogManager, portForwardController, watchManager, buildController, configsController, dockerComposeEventWatcher, dockerComposeLogManager, profilerManager, syncletManager, analyticsReporter, headsUpServerController, tiltVersionChecker, analyticsUpdater, eventWatchManager, cloudUsernameManager, updateUploader, dockerPruner)
	upper := engine.NewUpper(ctx, storeStore, v2)
	windmillDir, err := dirs.UseWindmillDir()
	if err != nil {
		return CmdUpDeps{}, err
	}
	tokenToken, err := token.GetOrCreateToken(windmillDir)
	if err != nil {
		return CmdUpDeps{}, err
	}
	cmdUpDeps := provideCmdUpDeps(headsUpDisplay, upper, tiltBuild, tokenToken, address)
	return cmdUpDeps, nil
}

var (
	_wireReducerValue = engine.UpperReducer
	_wireLabelsValue  = dockerfile.Labels{}
)

func wireKubeContext(ctx context.Context) (k8s.KubeContext, error) {
	clientConfig := k8s.ProvideClientConfig()
	config, err := k8s.ProvideKubeConfig(clientConfig)
	if err != nil {
		return "", err
	}
	kubeContext, err := k8s.ProvideKubeContext(config)
	if err != nil {
		return "", err
	}
	return kubeContext, nil
}

func wireKubeConfig(ctx context.Context) (*api.Config, error) {
	clientConfig := k8s.ProvideClientConfig()
	config, err := k8s.ProvideKubeConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func wireEnv(ctx context.Context) (k8s.Env, error) {
	clientConfig := k8s.ProvideClientConfig()
	config, err := k8s.ProvideKubeConfig(clientConfig)
	if err != nil {
		return "", err
	}
	env := k8s.ProvideEnv(ctx, config)
	return env, nil
}

func wireNamespace(ctx context.Context) (k8s.Namespace, error) {
	clientConfig := k8s.ProvideClientConfig()
	namespace := k8s.ProvideConfigNamespace(clientConfig)
	return namespace, nil
}

func wireClusterName(ctx context.Context) (k8s.ClusterName, error) {
	clientConfig := k8s.ProvideClientConfig()
	config, err := k8s.ProvideKubeConfig(clientConfig)
	if err != nil {
		return "", err
	}
	clusterName := k8s.ProvideClusterName(ctx, config)
	return clusterName, nil
}

func wireRuntime(ctx context.Context) (container.Runtime, error) {
	clientConfig := k8s.ProvideClientConfig()
	config, err := k8s.ProvideKubeConfig(clientConfig)
	if err != nil {
		return "", err
	}
	env := k8s.ProvideEnv(ctx, config)
	restConfigOrError := k8s.ProvideRESTConfig(clientConfig)
	clientsetOrError := k8s.ProvideClientset(restConfigOrError)
	portForwardClient := k8s.ProvidePortForwardClient(restConfigOrError, clientsetOrError)
	namespace := k8s.ProvideConfigNamespace(clientConfig)
	kubeContext, err := k8s.ProvideKubeContext(config)
	if err != nil {
		return "", err
	}
	int2 := provideKubectlLogLevel()
	kubectlRunner := k8s.ProvideKubectlRunner(kubeContext, int2)
	client := k8s.ProvideK8sClient(ctx, env, restConfigOrError, clientsetOrError, portForwardClient, namespace, kubectlRunner, clientConfig)
	runtime := k8s.ProvideContainerRuntime(ctx, client)
	return runtime, nil
}

func wireK8sVersion(ctx context.Context) (*version.Info, error) {
	clientConfig := k8s.ProvideClientConfig()
	restConfigOrError := k8s.ProvideRESTConfig(clientConfig)
	clientsetOrError := k8s.ProvideClientset(restConfigOrError)
	info, err := k8s.ProvideServerVersion(clientsetOrError)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func wireDockerClusterClient(ctx context.Context) (docker.ClusterClient, error) {
	clientConfig := k8s.ProvideClientConfig()
	config, err := k8s.ProvideKubeConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	env := k8s.ProvideEnv(ctx, config)
	restConfigOrError := k8s.ProvideRESTConfig(clientConfig)
	clientsetOrError := k8s.ProvideClientset(restConfigOrError)
	portForwardClient := k8s.ProvidePortForwardClient(restConfigOrError, clientsetOrError)
	namespace := k8s.ProvideConfigNamespace(clientConfig)
	kubeContext, err := k8s.ProvideKubeContext(config)
	if err != nil {
		return nil, err
	}
	int2 := provideKubectlLogLevel()
	kubectlRunner := k8s.ProvideKubectlRunner(kubeContext, int2)
	client := k8s.ProvideK8sClient(ctx, env, restConfigOrError, clientsetOrError, portForwardClient, namespace, kubectlRunner, clientConfig)
	runtime := k8s.ProvideContainerRuntime(ctx, client)
	minikubeClient := minikube.ProvideMinikubeClient()
	clusterEnv, err := docker.ProvideClusterEnv(ctx, env, runtime, minikubeClient)
	if err != nil {
		return nil, err
	}
	localEnv, err := docker.ProvideLocalEnv(ctx, clusterEnv)
	if err != nil {
		return nil, err
	}
	localClient := docker.ProvideLocalCli(ctx, localEnv)
	clusterClient, err := docker.ProvideClusterCli(ctx, localEnv, clusterEnv, localClient)
	if err != nil {
		return nil, err
	}
	return clusterClient, nil
}

func wireDockerLocalClient(ctx context.Context) (docker.LocalClient, error) {
	clientConfig := k8s.ProvideClientConfig()
	config, err := k8s.ProvideKubeConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	env := k8s.ProvideEnv(ctx, config)
	restConfigOrError := k8s.ProvideRESTConfig(clientConfig)
	clientsetOrError := k8s.ProvideClientset(restConfigOrError)
	portForwardClient := k8s.ProvidePortForwardClient(restConfigOrError, clientsetOrError)
	namespace := k8s.ProvideConfigNamespace(clientConfig)
	kubeContext, err := k8s.ProvideKubeContext(config)
	if err != nil {
		return nil, err
	}
	int2 := provideKubectlLogLevel()
	kubectlRunner := k8s.ProvideKubectlRunner(kubeContext, int2)
	client := k8s.ProvideK8sClient(ctx, env, restConfigOrError, clientsetOrError, portForwardClient, namespace, kubectlRunner, clientConfig)
	runtime := k8s.ProvideContainerRuntime(ctx, client)
	minikubeClient := minikube.ProvideMinikubeClient()
	clusterEnv, err := docker.ProvideClusterEnv(ctx, env, runtime, minikubeClient)
	if err != nil {
		return nil, err
	}
	localEnv, err := docker.ProvideLocalEnv(ctx, clusterEnv)
	if err != nil {
		return nil, err
	}
	localClient := docker.ProvideLocalCli(ctx, localEnv)
	return localClient, nil
}

func wireDownDeps(ctx context.Context, tiltAnalytics *analytics.TiltAnalytics) (DownDeps, error) {
	clientConfig := k8s.ProvideClientConfig()
	config, err := k8s.ProvideKubeConfig(clientConfig)
	if err != nil {
		return DownDeps{}, err
	}
	env := k8s.ProvideEnv(ctx, config)
	restConfigOrError := k8s.ProvideRESTConfig(clientConfig)
	clientsetOrError := k8s.ProvideClientset(restConfigOrError)
	portForwardClient := k8s.ProvidePortForwardClient(restConfigOrError, clientsetOrError)
	namespace := k8s.ProvideConfigNamespace(clientConfig)
	kubeContext, err := k8s.ProvideKubeContext(config)
	if err != nil {
		return DownDeps{}, err
	}
	int2 := provideKubectlLogLevel()
	kubectlRunner := k8s.ProvideKubectlRunner(kubeContext, int2)
	client := k8s.ProvideK8sClient(ctx, env, restConfigOrError, clientsetOrError, portForwardClient, namespace, kubectlRunner, clientConfig)
	extension := k8scontext.NewExtension(kubeContext, env)
	runtime := k8s.ProvideContainerRuntime(ctx, client)
	minikubeClient := minikube.ProvideMinikubeClient()
	clusterEnv, err := docker.ProvideClusterEnv(ctx, env, runtime, minikubeClient)
	if err != nil {
		return DownDeps{}, err
	}
	localEnv, err := docker.ProvideLocalEnv(ctx, clusterEnv)
	if err != nil {
		return DownDeps{}, err
	}
	dockerComposeClient := dockercompose.NewDockerComposeClient(localEnv)
	defaults := _wireDefaultsValue
	tiltfileLoader := tiltfile.ProvideTiltfileLoader(tiltAnalytics, client, extension, dockerComposeClient, defaults)
	downDeps := ProvideDownDeps(tiltfileLoader, dockerComposeClient, client)
	return downDeps, nil
}

// wire.go:

var K8sWireSet = wire.NewSet(k8s.ProvideEnv, k8s.DetectNodeIP, k8s.ProvideClusterName, k8s.ProvideKubeContext, k8s.ProvideKubeConfig, k8s.ProvideClientConfig, k8s.ProvideClientset, k8s.ProvideRESTConfig, k8s.ProvidePortForwardClient, k8s.ProvideConfigNamespace, k8s.ProvideKubectlRunner, k8s.ProvideContainerRuntime, k8s.ProvideServerVersion, k8s.ProvideK8sClient, k8s.ProvideOwnerFetcher)

var BaseWireSet = wire.NewSet(
	K8sWireSet, tiltfile.WireSet, provideKubectlLogLevel, docker.SwitchWireSet, dockercompose.NewDockerComposeClient, clockwork.NewRealClock, engine.DeployerWireSet, runtimelog.NewPodLogManager, engine.NewPortForwardController, engine.NewBuildController, k8swatch.NewPodWatcher, k8swatch.NewServiceWatcher, k8swatch.NewEventWatchManager, configs.NewConfigsController, engine.NewDockerComposeEventWatcher, runtimelog.NewDockerComposeLogManager, engine.NewProfilerManager, engine.NewGithubClientFactory, engine.NewTiltVersionChecker, cloud.WireSet, cloudurl.ProvideAddress, provideClock, hud.NewRenderer, hud.ProvideHud, provideLogActions, store.NewStore, wire.Bind(new(store.RStore), new(*store.Store)), dockerprune.NewDockerPruner, provideTiltInfo, engine.ProvideSubscribers, engine.NewUpper, analytics2.NewAnalyticsUpdater, analytics2.ProvideAnalyticsReporter, provideUpdateModeFlag, engine.NewWatchManager, engine.ProvideFsWatcherMaker, engine.ProvideTimerMaker, provideWebVersion,
	provideWebMode,
	provideWebURL,
	provideWebPort,
	provideWebHost,
	provideNoBrowserFlag, server.ProvideHeadsUpServer, provideAssetServer, server.ProvideHeadsUpServerController, dirs.UseWindmillDir, token.GetOrCreateToken, provideCmdUpDeps, engine.NewKINDPusher, wire.Value(feature.MainDefaults),
)

type CmdUpDeps struct {
	hud          hud.HeadsUpDisplay
	upper        engine.Upper
	tiltBuild    model.TiltBuild
	token        token.Token
	cloudAddress cloudurl.Address
}

func provideCmdUpDeps(h hud.HeadsUpDisplay, upper engine.Upper, b model.TiltBuild, token2 token.Token, cloudAddress cloudurl.Address) CmdUpDeps {
	return CmdUpDeps{h, upper, b, token2, cloudAddress}
}

type DownDeps struct {
	tfl      tiltfile.TiltfileLoader
	dcClient dockercompose.DockerComposeClient
	kClient  k8s.Client
}

func ProvideDownDeps(
	tfl tiltfile.TiltfileLoader,
	dcClient dockercompose.DockerComposeClient,
	kClient k8s.Client) DownDeps {
	return DownDeps{
		tfl:      tfl,
		dcClient: dcClient,
		kClient:  kClient,
	}
}

func provideClock() func() time.Time {
	return time.Now
}
