package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jaypipes/ghw"
	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/urfave/cli/v2"
)

func Status(cCtx *cli.Context) error {
	displayAll := cCtx.Bool("all")
	displayMemory := cCtx.Bool("memory")
	displayCPU := cCtx.Bool("cpu")
	displayGPU := cCtx.Bool("gpu")
	displayStorage := cCtx.Bool("storage")
	displayPlatform := cCtx.Bool("platform")

	if !displayMemory && !displayCPU && !displayGPU && !displayStorage && !displayPlatform {
		displayAll = true
	}

	if displayAll {
		displayMemory = true
		displayCPU = true
		displayGPU = true
		displayStorage = true
		displayPlatform = true
	}

	// Memory
	if displayMemory {
		color.Yellow("\nMemory")
		memory, err := ghw.Memory()
		if err != nil {
			return fmt.Errorf("error getting memory info: %v", err)
		}
		vmStat, _ := mem.VirtualMemory()
		memTable := tablewriter.NewWriter(os.Stdout)
		memTable.SetHeader([]string{"Total (MB)", "Used (MB)", "Free (MB)"})
		memTable.Append([]string{
			fmt.Sprintf("%v", vmStat.Total/1024/1024),
			fmt.Sprintf("%v", vmStat.Used/1024/1024),
			fmt.Sprintf("%v", vmStat.Free/1024/1024),
		})
		memTable.Render()

		fmt.Println(memory.String())
	}

	// CPU
	if displayCPU {
		color.Yellow("\nCPU")
		cpughw, err := ghw.CPU()
		if err != nil {
			return fmt.Errorf("error getting cpu info: %v", err)
		}
		cpuStat, err := cpu.Info()
		if err != nil {
			return fmt.Errorf("error getting cpu info: %v", err)
		}
		cpuTable := tablewriter.NewWriter(os.Stdout)
		cpuTable.SetHeader([]string{"Model", "Vendor", "Cache (KB)", "Cores", "Threads", "Mhz"})
		cpuTable.Append([]string{
			cpuStat[0].ModelName,
			cpuStat[0].VendorID,
			fmt.Sprintf("%v", cpuStat[0].CacheSize),
			fmt.Sprintf("%v", cpughw.Processors[0].NumCores),
			fmt.Sprintf("%v", cpughw.Processors[0].NumThreads),
			fmt.Sprintf("%v", cpuStat[0].Mhz),
		})
		cpuTable.Render()

		for _, proc := range cpughw.Processors {
			fmt.Printf(" Processor: %v\n", proc)
			for _, core := range proc.Cores {
				fmt.Printf("  Core: %v\n", core)
			}
			if len(proc.Capabilities) > 0 {
				capStr := strings.Join(proc.Capabilities, ", ")
				fmt.Printf("  Capabilities: %s\n", capStr)
			}
		}
	}

	// GPU
	if displayGPU {
		color.Yellow("\nGPU")
		gpu, err := ghw.GPU()
		if err != nil {
			return fmt.Errorf("error getting gpu info: %v", err)
		}
		gpuTable := tablewriter.NewWriter(os.Stdout)
		gpuTable.SetHeader([]string{"Index", "Vendor", "Product"})
		for _, card := range gpu.GraphicsCards {
			gpuTable.Append([]string{
				fmt.Sprintf("%v", card.Index),
				card.DeviceInfo.Vendor.Name,
				card.DeviceInfo.Product.Name,
			})
		}
		gpuTable.Render()
	}

	// Platform
	if displayPlatform {
		color.Yellow("\nPlatform")
		hostStat, err := host.Info()
		if err != nil {
			return fmt.Errorf("error getting platform info: %v", err)
		}
		platformTable := tablewriter.NewWriter(os.Stdout)
		platformTable.SetHeader([]string{"Hostname", "Platform", "OS", "Uptime (min)"})
		platformTable.Append([]string{
			hostStat.Hostname,
			hostStat.Platform,
			hostStat.OS,
			fmt.Sprintf("%v", hostStat.Uptime/60),
		})
		platformTable.Render()
	}

	// Storage
	if displayStorage {
		color.Yellow("\nStorage")
		block, err := ghw.Block()
		if err != nil {
			return fmt.Errorf("error getting block storage info: %v", err)
		}
		storageTable := tablewriter.NewWriter(os.Stdout)
		storageTable.SetHeader([]string{"Name", "Size (GB)", "Vendor", "Partitions"})

		for _, disk := range block.Disks {
			var partitions []string
			for _, part := range disk.Partitions {
				partitions = append(partitions, part.Name)
			}
			storageTable.Append([]string{
				disk.Name,
				fmt.Sprintf("%.2f", float64(disk.SizeBytes)/1e9),
				disk.Vendor,
				strings.Join(partitions, ", "),
			})
		}
		storageTable.Render()
	}
	return nil
}
