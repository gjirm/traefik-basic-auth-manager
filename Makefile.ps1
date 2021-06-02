#!/usr/bin/pwsh

$projectName = "Traefik Basic Auth Manager"
#$dt = Get-Date -Format "yyy-MM-dd_HHMMss"
$headhash = git rev-parse --short HEAD
#$tag = "$($dt)-$($headhash)"
$tag = git describe --tags --abbrev=0
$imageName = "jirm/tbam:$($headhash)"
$imageLatest = "jirm/tbam:latest"
$imageNameTag = "jirm/tbam:$($tag)"
$shortName = "tbam"

#$minisignKey = "W:\keys\jirm-minisign-2020.key"

#docker build --no-cache --tag $(IMAGENAMETAG) --tag $(IMAGENAME) --tag $(IMAGELATEST) .

Write-Host "--> $projectName <--" -ForegroundColor Green
if ($Args[0] -eq "build-docker-nocache-tag") {
    
    Write-Host "--> Building $($imageNameTag)" -ForegroundColor Green
    docker build --no-cache --tag $imageNameTag --tag $imageName --tag $imageLatest .
    If ($lastExitCode -eq "0") {
        Write-Host "--> $($imageName) successfully build!" -ForegroundColor Green
        exit 0
    } else {
        Write-Host "--X $($imageName) build failed!" -ForegroundColor Red
        exit 1
    }
}

if ($Args[0] -eq "build-docker-tag") {
    
    Write-Host "--> Building $($imageNameTag)" -ForegroundColor Green
    docker build --tag $imageNameTag --tag $imageName --tag $imageLatest .
    If ($lastExitCode -eq "0") {
        Write-Host "--> $($imageName) successfully build!" -ForegroundColor Green
        exit 0
    } else {
        Write-Host "--X $($imageName) build failed!" -ForegroundColor Red
        exit 1
    }
}

if ($Args[0] -eq "build-docker-nonametag") {
    
    Write-Host "--> Building $($imageNameTag)" -ForegroundColor Green
    docker build --tag $imageName --tag $imageLatest .
    If ($lastExitCode -eq "0") {
        Write-Host "--> $($imageName) successfully build!" -ForegroundColor Green
        exit 0
    } else {
        Write-Host "--X $($imageName) build failed!" -ForegroundColor Red
        exit 1
    }
} 

if ($args[0] -eq "run") {

    Write-Host "--> Running $shortName..."  -ForegroundColor Green
    go run .\app\main.go
    exit 0

}

if ($args[0] -eq "run-docker") {

        Write-Host "--> Run Docker container"  -ForegroundColor Green
        try {
            docker run --rm -v $PSScriptRoot/config.yml:/tbam/config.yml -v $PSScriptRoot/db:/tbam/nutsdb -v $PSScriptRoot/basic-auth.yml:/tbam/basic-auth.yml --name $shortName -p 8080:8080 $imageLatest
        }
        catch {
            Write-Error $_.Exception
        }
        finally {
            docker stop $shortName
        }
        exit 0

}

Write-Host "--! None!" -ForegroundColor Yellow
