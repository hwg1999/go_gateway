package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/hwg1999/go_gateway/backend/dao"
	"github.com/hwg1999/go_gateway/backend/dto"
	"github.com/hwg1999/go_gateway/backend/golang_common/lib"
	"github.com/hwg1999/go_gateway/backend/middleware"
	"github.com/hwg1999/go_gateway/backend/public"
)

type DashboardController struct{}

func DashboardRegister(group *gin.RouterGroup) {
	service := &DashboardController{}
	group.GET("/panel_group_data", service.PanelGroupData)
	group.GET("/service_stat", service.ServiceStat)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]
func (service *DashboardController) PanelGroupData(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{}
	_, serviceNum, err := serviceInfo.PageList(c, tx, &dto.ServiceListInput{PageSize: 1, PageNo: 1})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	app := &dao.App{}
	_, appNum, err := app.APPList(c, tx, &dto.APPListInput{PageNo: 1, PageSize: 1})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	out := &dto.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		AppNum:          appNum,
		TodayRequestNum: 0,
		CurrentQPS:      0,
	}

	middleware.ResponseSuccess(c, out)
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.DashServiceStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (service *DashboardController) ServiceStat(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	list, err := serviceInfo.GroupByLoadType(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	legend := []string{}
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(c, 2003, errors.New("load_type not found"))
			return
		}
		list[index].Name = name
		legend = append(legend, name)
	}
	out := &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	middleware.ResponseSuccess(c, out)
}
